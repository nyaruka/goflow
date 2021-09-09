package wit

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

const (
	apiBaseURL = "https://api.wit.ai"
	version    = "20200513"
)

// IntentMatch is possible intent match
type IntentMatch struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Confidence decimal.Decimal `json:"confidence"`
}

// EntityMatch is possible entity match
type EntityMatch struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Role       string          `json:"role"`
	Value      string          `json:"value"`
	Confidence decimal.Decimal `json:"confidence"`
}

// TraitMatch is possible trait match
type TraitMatch struct {
	ID         string          `json:"id"`
	Value      string          `json:"value"`
	Confidence decimal.Decimal `json:"confidence"`
}

// MessageResponse is the response from a /message request
type MessageResponse struct {
	Text     string                   `json:"text"`
	Intents  []IntentMatch            `json:"intents" validate:"required"`
	Entities map[string][]EntityMatch `json:"entities"`
	Traits   map[string][]TraitMatch  `json:"traits"`
}

// Client is a basic Wit.ai client
type Client struct {
	httpClient  *http.Client
	httpRetries *httpx.RetryConfig
	headers     map[string]string
}

// NewClient creates a new client
func NewClient(httpClient *http.Client, httpRetries *httpx.RetryConfig, accessToken string) *Client {
	return &Client{
		httpClient:  httpClient,
		httpRetries: httpRetries,
		headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		},
	}
}

// Message gets the meaning of a message
func (c *Client) Message(q string) (*MessageResponse, *httpx.Trace, error) {
	endpoint := fmt.Sprintf("%s/message?v=%s&q=%s", apiBaseURL, version, url.QueryEscape(q))

	request, err := httpx.NewRequest("GET", endpoint, nil, c.headers)
	if err != nil {
		return nil, nil, err
	}

	trace, err := httpx.DoTrace(c.httpClient, request, c.httpRetries, nil, -1)
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

	return nil, trace, errors.New("wit.ai API request failed")
}
