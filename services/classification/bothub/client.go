package bothub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

const (
	apiBaseURL = "https://nlp.bothub.it"
)

type IntentMatch struct {
	Name       string          `json:"name"`
	Confidence decimal.Decimal `json:"confidence"`
}

type EntityMatch struct {
	Value      string          `json:"value"`
	Entity     string          `json:"entity"`
	Confidence decimal.Decimal `json:"confidence"`
}

// ParseResponse is the response from a /parse request
type ParseResponse struct {
	Intent        IntentMatch              `json:"intent" validate:"required"`
	IntentRanking []IntentMatch            `json:"intent_ranking" validate:"required"`
	LabelsList    []string                 `json:"labels_list"`
	EntitiesList  []string                 `json:"entities_list"`
	Entities      map[string][]EntityMatch `json:"entities"`
	Text          string                   `json:"text"`
	UpdateID      int                      `json:"update_id"`
	Language      string                   `json:"language"`
}

// Client is a basic Bothub client
type Client struct {
	httpClient  *http.Client
	httpRetries *httpx.RetryConfig
	accessToken string
}

// NewClient creates a new client
func NewClient(httpClient *http.Client, httpRetries *httpx.RetryConfig, accessToken string) *Client {
	return &Client{
		httpClient:  httpClient,
		httpRetries: httpRetries,
		accessToken: accessToken,
	}
}

// Parse does a parse of the given text in the given language (e.g. pt_br)
func (c *Client) Parse(text, language string) (*ParseResponse, *httpx.Trace, error) {
	endpoint := fmt.Sprintf("%s/parse", apiBaseURL)

	form := url.Values{}
	form.Add("text", text)
	if language != "" {
		form.Add("language", language)
	}

	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": fmt.Sprintf("Bearer %s", c.accessToken),
	}

	request, err := httpx.NewRequest("POST", endpoint, strings.NewReader(form.Encode()), headers)
	if err != nil {
		return nil, nil, err
	}

	trace, err := httpx.DoTrace(c.httpClient, request, c.httpRetries, nil, -1)
	if err != nil {
		return nil, trace, err
	}

	response := &ParseResponse{}
	if err := utils.UnmarshalAndValidate(trace.ResponseBody, response); err != nil {
		return nil, trace, err
	}

	return response, trace, nil
}
