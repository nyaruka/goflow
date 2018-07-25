package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

const (
	apiURL = "https://airtime.transferto.com/cgi-bin/shop/topup"
)

// Client is a TransferTo client
// see https://shop.transferto.com/shop/v3/doc/TransferTo_API.pdf for API docs
type Client struct {
	login      string
	token      string
	httpClient *utils.HTTPClient
	apiURL     string
}

// Response is a base interface for all responses
type Response interface {
	ErrorCode() int
	ErrorTxt() string
}

type baseResponse struct {
	ErrorCode_ int    `json:"error_code,string"`
	ErrorTxt_  string `json:"error_txt"`
}

func (r *baseResponse) ErrorCode() int   { return r.ErrorCode_ }
func (r *baseResponse) ErrorTxt() string { return r.ErrorTxt_ }

// NewTransferToClient creates a new TransferTo client
func NewTransferToClient(login string, token string, httpClient *utils.HTTPClient) *Client {
	return &Client{login: login, token: token, httpClient: httpClient, apiURL: apiURL}
}

// SetAPIURL sets the API URL used by this client
func (c *Client) SetAPIURL(url string) {
	c.apiURL = url
}

// Ping just verifies the credentials
func (c *Client) Ping() error {
	request := url.Values{}
	request.Add("action", "ping")

	response := &struct {
		baseResponse
		InfoTxt string `json:"info_txt"`
	}{}
	return c.request(request, response)
}

// MSISDNInfo is a response to a msisdn_info request
type MSISDNInfo struct {
	baseResponse
	Country             string      `json:"country"`
	CountryID           int         `json:"country_id,string"`
	Operator            string      `json:"operator"`
	OperatorID          int         `json:"operator_id,string"`
	ConnectionStatus    int         `json:"connection_status,string"`
	DestinationCurrency string      `json:"destination_currency"`
	ProductList         CSVStrings  `json:"product_list"`
	ServiceFeeList      CSVDecimals `json:"service_fee_list"`
	SKUIDList           CSVStrings  `json:"skuid_list"`
	LocalInfoValueList  CSVDecimals `json:"local_info_value_list"`

	// if operator supports open-range transfers...
	OpenRange                           bool            `json:"open_range"`
	SKUID                               string          `json:"skuid"`
	OpenRangeMinimumAmountLocalCurrency decimal.Decimal `json:"open_range_minimum_amount_local_currency"`
	OpenRangeMaximumAmountLocalCurrency decimal.Decimal `json:"open_range_maximum_amount_local_currency"`
	OpenRangeIncrementLocalCurrency     decimal.Decimal `json:"open_range_increment_local_currency"`
}

// MSISDNInfo fetches information about the given MSISDN
func (c *Client) MSISDNInfo(destinationMSISDN string, currency string, deliveredAmountInfo string) (*MSISDNInfo, error) {
	request := url.Values{}
	request.Add("action", "msisdn_info")
	request.Add("destination_msisdn", destinationMSISDN)
	request.Add("currency", currency)
	request.Add("delivered_amount_info", deliveredAmountInfo)

	response := &MSISDNInfo{}
	if err := c.request(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// ReserveID reserves a transaction ID for a future topup
func (c *Client) ReserveID() (int, error) {
	request := url.Values{}
	request.Add("action", "reserve_id")

	response := &struct {
		baseResponse
		ReservedID int `json:"reserved_id,string"`
	}{}
	if err := c.request(request, response); err != nil {
		return 0, err
	}
	return response.ReservedID, nil
}

// Topup is a response to a topup request
type Topup struct {
	baseResponse
	DestinationCurrency string          `json:"destination_currency"`
	OriginatingCurrency string          `json:"originating_currency"`
	ProductRequested    decimal.Decimal `json:"product_requested"`
	ActualProductSent   decimal.Decimal `json:"actual_product_sent"`
	SMSSent             string          `json:"sms_sent"`
	SMS                 string          `json:"sms"`
	WholesalePrice      decimal.Decimal `json:"wholesale_price"`
	ServiceFee          decimal.Decimal `json:"service_fee"`
	RetailPrice         decimal.Decimal `json:"retail_price"`
	LocalInfoAmount     decimal.Decimal `json:"local_info_amount"`
	LocalInfoValue      decimal.Decimal `json:"local_info_value"`
	Balance             decimal.Decimal `json:"balance"`
}

// Topup makes an actual airtime transfer
func (c *Client) Topup(reservedID int, msisdn string, destinationMSISDN string, product string, skuid string) (*Topup, error) {
	request := url.Values{}
	request.Add("action", "topup")
	request.Add("reserved_id", strconv.Itoa(reservedID))
	request.Add("msisdn", msisdn)
	request.Add("destination_msisdn", destinationMSISDN)
	request.Add("product", product)
	if skuid != "" {
		request.Add("skuid", skuid)
	}

	response := &Topup{}
	if err := c.request(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// makes a request with the given data and parses the response into the destination struct
func (c *Client) request(data url.Values, dest interface{}) error {
	key := strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond)))

	hasher := md5.New()
	hasher.Write([]byte(c.login + c.token + key))
	hash := hex.EncodeToString(hasher.Sum(nil))

	data.Add("login", c.login)
	data.Add("key", key)
	data.Add("md5", hash)

	req, err := http.NewRequest("POST", c.apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, _, err := c.httpClient.DoWithDump(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if err := c.parseResponse(response.Body, dest); err != nil {
		return fmt.Errorf("transferto API call returned an unparseable response: %s", err)
	}

	baseResp := dest.(Response)
	if baseResp.ErrorCode() != 0 {
		return fmt.Errorf("transferto API call returned an error: %s (%d)", baseResp.ErrorTxt(), baseResp.ErrorCode())
	}
	return nil
}

// reads and parses a response body, which is in the format
//
// value1=result1
// value2=result2
// ...
//
// with each line separated by \r\n
func (c *Client) parseResponse(body io.Reader, dest interface{}) error {
	asBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	// parse into a map
	data := make(map[string]string)
	for _, line := range strings.Split(string(asBytes), "\r\n") {
		parts := strings.SplitN(line, "=", 2)

		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}

	// marshal to JSON and then into the destination struct
	respJSON, _ := json.Marshal(data)
	return json.Unmarshal(respJSON, dest)
}
