package dtone

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/shopspring/decimal"
)

const (
	apiURL = "https://dvs-api.dtone.com/v1/"
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
	for m := range e.Errors {
		msgs[m] = e.Errors[m].Message
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
func (c *Client) LookupMobileNumber(tel string) ([]*Operator, *httpx.Trace, error) {
	var response []*Operator

	trace, err := c.request("GET", fmt.Sprintf("lookup/mobile-number/%s", tel), nil, &response)
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
func (c *Client) Products(_type string, operatorID int) ([]*Product, *httpx.Trace, error) {
	var response []*Product

	trace, err := c.request("GET", fmt.Sprintf("products?type=%s&operator_id=%d", _type, operatorID), nil, &response)
	if err != nil {
		return nil, trace, err
	}

	return response, trace, nil
}

func (c *Client) request(method, endpoint string, payload interface{}, response interface{}) (*httpx.Trace, error) {
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

	req, err := httpx.NewRequest(method, url, body, headers)
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
