package dtone

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"

	"github.com/shopspring/decimal"
)

type StatusCID int

const (
	apiURL = "https://dvs-api.dtone.com/v1/"

	// see https://dvs-api-doc.dtone.com/#section/Overview/Transactions
	StatusCIDCreated   StatusCID = 1
	StatusCIDConfirmed StatusCID = 2
	StatusCIDRejected  StatusCID = 3
	StatusCIDCancelled StatusCID = 4
	StatusCIDSubmitted StatusCID = 5
	StatusCIDCompleted StatusCID = 7
	StatusCIDReversed  StatusCID = 8
	StatusCIDDeclined  StatusCID = 9
)

// Client is a DTOne client, see https://dvs-api-doc.dtone.com/ for API docs
type Client struct {
	httpClient  *http.Client
	httpRetries *httpx.RetryConfig
	key         string
	secret      string
}

// NewClient creates a new DT One client
func NewClient(httpClient *http.Client, httpRetries *httpx.RetryConfig, key, secret string) *Client {
	return &Client{httpClient: httpClient, httpRetries: httpRetries, key: key, secret: secret}
}

// error response contains errors when a request fails
type errorResponse struct {
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func (e *errorResponse) Error() string {
	msgs := make([]string, len(e.Errors))
	for i := range e.Errors {
		msgs[i] = e.Errors[i].Message
	}
	return strings.Join(msgs, ", ")
}

// Operator is a mobile operator
type Operator struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country struct {
		Name    string `json:"name"`
		ISOCode string `json:"iso_code"`
		Regions []struct {
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"regions"`
	} `json:"country"`
	Identified bool `json:"identified"`
}

// LookupMobileNumber see https://dvs-api-doc.dtone.com/#tag/Mobile-Number
func (c *Client) LookupMobileNumber(ctx context.Context, phoneNumber string) ([]*Operator, *httpx.Trace, error) {
	var response []*Operator

	payload := &struct {
		MobileNumber string `json:"mobile_number"`
	}{
		MobileNumber: phoneNumber,
	}

	trace, err := c.request(ctx, "POST", "lookup/mobile-number", payload, &response)
	if err != nil {
		return nil, trace, err
	}

	return response, trace, nil
}

// Product is an available digital services product
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Service     struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"service"`
	Operator struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"operator"`
	Type   string `json:"type"`
	Source struct {
		Amount   decimal.Decimal `json:"amount"`
		Unit     string          `json:"unit"`
		UnitType string          `json:"unit_type"`
	} `json:"source"`
	Destination struct {
		Amount   decimal.Decimal `json:"amount"`
		Unit     string          `json:"unit"`
		UnitType string          `json:"unit_type"`
	} `json:"destination"`
}

// Products see https://dvs-api-doc.dtone.com/#tag/Products
func (c *Client) Products(ctx context.Context, _type string, operatorID int) ([]*Product, *httpx.Trace, error) {
	var response []*Product

	// TODO endpoint could return more than 100 products in which case we need to page

	trace, err := c.request(ctx, "GET", fmt.Sprintf("products?type=%s&operator_id=%d&per_page=100", _type, operatorID), nil, &response)
	if err != nil {
		return nil, trace, err
	}

	return response, trace, nil
}

// Transaction is a product sent to a beneficiary
type Transaction struct {
	ID                         int64  `json:"id"`
	ExternalID                 string `json:"external_id"`
	CreationDate               string `json:"creation_date"`
	ConfirmationExpirationDate string `json:"confirmation_expiration_date"`
	ConfirmationDate           string `json:"confirmation_date"`
	Status                     struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
		Class   struct {
			ID      StatusCID `json:"id"`
			Message string    `json:"message"`
		}
	} `json:"status"`
}

// TransactionAsync see https://dvs-api-doc.dtone.com/#tag/Transactions
func (c *Client) TransactionAsync(ctx context.Context, externalID string, productID int, mobileNumber string) (*Transaction, *httpx.Trace, error) {
	var response *Transaction

	type creditPartyIdentifier struct {
		MobileNumber string `json:"mobile_number"`
	}

	payload := &struct {
		ExternalID            string                `json:"external_id"`
		ProductID             int                   `json:"product_id"`
		AutoConfirm           bool                  `json:"auto_confirm"`
		CreditPartyIdentifier creditPartyIdentifier `json:"credit_party_identifier"`
	}{
		ExternalID:  externalID,
		ProductID:   productID,
		AutoConfirm: true,
		CreditPartyIdentifier: creditPartyIdentifier{
			MobileNumber: mobileNumber,
		},
	}

	trace, err := c.request(ctx, "POST", "async/transactions", payload, &response)
	if err != nil {
		return nil, trace, err
	}

	return response, trace, nil
}

func (c *Client) request(ctx context.Context, method, endpoint string, payload any, response any) (*httpx.Trace, error) {
	url := apiURL + endpoint
	headers := map[string]string{}
	var body io.Reader

	if payload != nil {
		data, err := jsonx.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
		headers["Content-Type"] = "application/json"
	}

	req, err := httpx.NewRequest(ctx, method, url, body, headers)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.key, c.secret)

	trace, err := httpx.DoTrace(c.httpClient, req, c.httpRetries, nil, -1)
	if err != nil {
		return trace, err
	}

	if trace.Response.StatusCode >= 400 {
		response := &errorResponse{}
		jsonx.Unmarshal(trace.ResponseBody, response)
		return trace, response
	}

	if response != nil {
		return trace, jsonx.Unmarshal(trace.ResponseBody, response)
	}
	return trace, nil
}
