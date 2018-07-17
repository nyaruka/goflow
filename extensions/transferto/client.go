package transferto

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"
)

const (
	apiURL = "https://airtime.transferto.com/cgi-bin/shop/topup"
)

type Client struct {
	login      string
	token      string
	httpClient *utils.HTTPClient
}

func NewTransferToClient(login string, token string, httpClient *utils.HTTPClient) *Client {
	return &Client{login: login, token: token, httpClient: httpClient}
}

func (c *Client) MSISDNInfo(msisdn string, currency string) (map[string]string, error) {
	request := url.Values{}
	request.Add("action", "msisdn_info")
	request.Add("destination_msisdn", msisdn)
	request.Add("currency", currency)
	request.Add("delivered_amount_info", "1")
	return c.request(request)
}

func (c *Client) ReserveID() (int, error) {
	request := url.Values{}
	request.Add("action", "reserve_id")

	response, err := c.request(request)
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(response["reserved_id"])
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *Client) Topup(reservedID int) (map[string]string, error) {
	data := url.Values{}
	data.Add("action", "topup")
	data.Add("reserved_id", strconv.Itoa(reservedID))
	return c.request(data)
}

func (c *Client) request(data url.Values) (map[string]string, error) {
	key := strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond)))

	hasher := md5.New()
	hasher.Write([]byte(c.login + c.token + key))
	hash := hex.EncodeToString(hasher.Sum(nil))

	data.Add("login", c.login)
	data.Add("key", key)
	data.Add("md5", hash)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("transferto API call return non-2XX response (%d)", response.StatusCode)
	}

	defer response.Body.Close()
	respData, err := c.parseResponse(response.Body)
	if err != nil {
		return nil, fmt.Errorf("transferto API call returned an unparseable response")
	}

	errorCode := respData["error_code"]
	errorText := respData["error_txt"]

	if errorCode != "" {
		return nil, fmt.Errorf("transferto API call returned an error: %s (%s)", errorText, errorCode)
	}

	return respData, nil
}

// reads and parses a response body, which is in the format
//
// value1=result1
// value2=result2
// ...
//
// with each line separated by \r\n
func (c *Client) parseResponse(body io.Reader) (map[string]string, error) {
	asBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	for _, line := range strings.Split(string(asBytes), "\r\n") {
		parts := strings.SplitN(line, "=", 2)
		data[parts[0]] = parts[1]
	}

	return data, nil
}
