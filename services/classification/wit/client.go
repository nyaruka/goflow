package wit

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/shopspring/decimal"
)

const (
	apiBaseURL = "https://api.wit.ai"
	version    = "20170307"
)

type EntityCandidate struct {
	Value      string          `json:"value"`
	Confidence decimal.Decimal `json:"confidence"`
}

// MessageResponse is the response from a /message request
type MessageResponse struct {
	MsgID    string                       `json:"msg_id"`
	Text     string                       `json:"_text"`
	Entities map[string][]EntityCandidate `json:"entities" validate:"required"`
}

// Client is a basic Wit.ai client
type Client struct {
	httpClient *http.Client
	headers    map[string]string
}

// NewClient creates a new client
func NewClient(httpClient *http.Client, accessToken string) *Client {
	return &Client{
		httpClient: httpClient,
		headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		},
	}
}

// Message gets the meaning of a message
func (c *Client) Message(q string) (*MessageResponse, *httpx.Trace, error) {
	endpoint := fmt.Sprintf("%s/message?v=%s&q=%s", apiBaseURL, version, url.QueryEscape(q))

	trace, err := httpx.DoTrace(c.httpClient, "GET", endpoint, nil, c.headers, nil)
	if err != nil {
		return nil, trace, err
	}

	if trace.Response != nil && trace.Response.StatusCode == 200 {
		response := &MessageResponse{}
		if err := utils.UnmarshalAndValidate(trace.ResponseBody, response); err != nil {
			return nil, trace, err
		}
		return response, trace, nil
	}

	return nil, trace, errors.New("Wit API request failed")
}
