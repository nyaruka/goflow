package transferto

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
type Client struct {
	login      string
	token      string
	httpClient *utils.HTTPClient
}

type baseResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorTxt  string `json:"error_txt"`
}

// MSISDNInfo holds information about an MSISDN
type MSISDNInfo struct {
	baseResponse
	Country             string  `json:"country"`
	CountryID           int     `json:"country_id"`
	Operator            string  `json:"operator"`
	OperatorID          int     `json:"operator_id"`
	ConnectionStatus    int     `json:"connection_status"`
	DestinationCurrency string  `json:"destination_currency"`
	ProductList         CSVList `json:"product_list"`
	ServiceFeeList      CSVList `json:"service_fee_list"`
	SKUIDList           CSVList `json:"skuid_list"`
	LocalInfoValueList  CSVList `json:"local_info_value_list"`
}

// NewTransferToClient creates a new TransferTo client
func NewTransferToClient(login string, token string, httpClient *utils.HTTPClient) *Client {
	return &Client{login: login, token: token, httpClient: httpClient}
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
		ReservedID int `json:"reserved_id"`
	}{}
	if err := c.request(request, response); err != nil {
		return 0, err
	}
	return response.ReservedID, nil
}

// Topup makes an actual airtime transfer
func (c *Client) Topup(reservedID int) (decimal.Decimal, error) {
	request := url.Values{}
	request.Add("action", "topup")
	request.Add("reserved_id", strconv.Itoa(reservedID))

	response := &struct {
		baseResponse
		Balance decimal.Decimal `json:"balance"`
	}{}
	if err := c.request(request, response); err != nil {
		return decimal.Zero, err
	}
	return response.Balance, nil
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

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	//if response.StatusCode < 200 || response.StatusCode >= 300 {
	//	return nil, fmt.Errorf("transferto API call return non-2XX response (%d)", response.StatusCode)
	//}

	defer response.Body.Close()
	if err := c.parseResponse(response.Body, dest); err != nil {
		return fmt.Errorf("transferto API call returned an unparseable response")
	}

	baseResp := dest.(*baseResponse)
	if baseResp.ErrorCode != 0 {
		return fmt.Errorf("transferto API call returned an error: %s (%d)", baseResp.ErrorTxt, baseResp.ErrorCode)
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
		data[parts[0]] = parts[1]
	}

	// marshal to JSON and then into the destination struct
	respJSON, _ := json.Marshal(data)
	return json.Unmarshal(respJSON, dest)
}
