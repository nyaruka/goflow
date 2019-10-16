package dtone

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	apiURL = "https://airtime-api.dtone.com/cgi-bin/shop/topup"
)

// Client is a TransferTo client
// see https://tshop-app.dtone.com/shop/v3/doc/Airtime_API.pdf for API docs
type Client struct {
	httpClient *http.Client
	login      string
	token      string
}

type BaseResponse struct {
	ErrorCode int    `json:"error_code,string"`
	ErrorTxt  string `json:"error_txt"`
}

// NewClient creates a new TransferTo client
func NewClient(httpClient *http.Client, login string, token string) *Client {
	return &Client{httpClient: httpClient, login: login, token: token}
}

// Ping just verifies the credentials
func (c *Client) Ping() (*httpx.Trace, error) {
	request := url.Values{}
	request.Add("action", "ping")

	response := &struct {
		BaseResponse
		InfoTxt string `json:"info_txt"`
	}{}
	return c.request(request, response)
}

// MSISDNInfo is a response to a msisdn_info request
type MSISDNInfo struct {
	BaseResponse
	Country             string      `json:"country"`
	CountryID           int         `json:"country_id,string"`
	Operator            string      `json:"operator"`
	OperatorID          int         `json:"operator_id,string"`
	ConnectionStatus    int         `json:"connection_status,string"`
	DestinationCurrency string      `json:"destination_currency" validate:"required"`
	ProductList         CSVStrings  `json:"product_list" validate:"required"`
	ServiceFeeList      CSVDecimals `json:"service_fee_list"`
	SKUIDList           CSVStrings  `json:"skuid_list"`
	LocalInfoValueList  CSVDecimals `json:"local_info_value_list" validate:"required"`

	// if operator supports open-range transfers...
	/*OpenRange                           bool            `json:"open_range"`
	SKUID                               string          `json:"skuid"`
	OpenRangeMinimumAmountLocalCurrency decimal.Decimal `json:"open_range_minimum_amount_local_currency"`
	OpenRangeMaximumAmountLocalCurrency decimal.Decimal `json:"open_range_maximum_amount_local_currency"`
	OpenRangeIncrementLocalCurrency     decimal.Decimal `json:"open_range_increment_local_currency"`*/
}

// MSISDNInfo fetches information about the given MSISDN
func (c *Client) MSISDNInfo(destinationMSISDN string, currency string, deliveredAmountInfo string) (*MSISDNInfo, *httpx.Trace, error) {
	request := url.Values{}
	request.Add("action", "msisdn_info")
	request.Add("destination_msisdn", destinationMSISDN)
	request.Add("currency", currency)
	request.Add("delivered_amount_info", deliveredAmountInfo)

	response := &MSISDNInfo{}
	trace, err := c.request(request, response)
	if err != nil {
		return nil, trace, err
	}
	return response, trace, nil
}

// ReserveID is a response to a reserve_id request
type ReserveID struct {
	BaseResponse
	ReservedID int `json:"reserved_id,string" validate:"required"`
}

// ReserveID reserves a transaction ID for a future topup
func (c *Client) ReserveID() (int, *httpx.Trace, error) {
	request := url.Values{}
	request.Add("action", "reserve_id")

	response := &ReserveID{}
	trace, err := c.request(request, response)
	if err != nil {
		return 0, trace, err
	}
	return response.ReservedID, trace, nil
}

// Topup is a response to a topup request
type Topup struct {
	BaseResponse
	DestinationCurrency string          `json:"destination_currency" validate:"required"`
	OriginatingCurrency string          `json:"originating_currency"`
	ProductRequested    decimal.Decimal `json:"product_requested"`
	ActualProductSent   decimal.Decimal `json:"actual_product_sent" validate:"required"`
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
func (c *Client) Topup(reservedID int, msisdn string, destinationMSISDN string, product string, skuid string) (*Topup, *httpx.Trace, error) {
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
	trace, err := c.request(request, response)
	if err != nil {
		return nil, trace, err
	}

	return response, trace, nil
}

// makes a request with the given data and parses the response into the destination struct
func (c *Client) request(data url.Values, dest interface{}) (*httpx.Trace, error) {
	key := strconv.Itoa(int(dates.Now().UnixNano() / int64(time.Millisecond)))

	hasher := md5.New()
	hasher.Write([]byte(c.login + c.token + key))
	hash := hex.EncodeToString(hasher.Sum(nil))

	data.Add("login", c.login)
	data.Add("key", key)
	data.Add("md5", hash)

	trace, err := httpx.DoTrace(c.httpClient, "POST", apiURL, strings.NewReader(data.Encode()), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	if err != nil {
		return nil, err
	}

	if err := c.parseResponse(trace.ResponseBody, dest); err != nil {
		return trace, errors.Wrap(err, "DTOne API request failed")
	}

	return trace, nil
}

// reads and parses a response body, which is in the format
//
// value1=result1
// value2=result2
// ...
//
// with each line separated by \r\n
func (c *Client) parseResponse(asBytes []byte, dest interface{}) error {
	// parse into a map
	data := make(map[string]string)
	for _, line := range strings.Split(string(asBytes), "\r\n") {
		parts := strings.SplitN(line, "=", 2)

		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}

	// marshal to JSON so we can use nice golang JSON unmarshalling into our response structs
	respJSON, _ := json.Marshal(data)

	// first try to unmarshal as base response which contains error messages
	baseResp := &BaseResponse{}
	utils.UnmarshalAndValidate(respJSON, baseResp)

	if baseResp.ErrorCode != 0 {
		return errors.Errorf("%s (%d)", baseResp.ErrorTxt, baseResp.ErrorCode)
	}

	// now unmarshal into action specific struct
	return utils.UnmarshalAndValidate(respJSON, dest)
}
