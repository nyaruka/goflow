package luis

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/shopspring/decimal"
)

type ExtractedIntent struct {
	Intent string          `json:"intent"`
	Score  decimal.Decimal `json:"score"`
}

type ExtractedEntity struct {
	Entity     string          `json:"entity"`
	Type       string          `json:"type"`
	StartIndex int             `json:"startIndex"`
	EndIndex   int             `json:"endIndex"`
	Score      decimal.Decimal `json:"score"`
}

type SentimentAnalysis struct {
	Label string          `json:"label"`
	Score decimal.Decimal `json:"score"`
}

// PredictResponse is the response from a predict request
type PredictResponse struct {
	Query             string             `json:"query"`
	TopScoringIntent  *ExtractedIntent   `json:"topScoringIntent"`
	Intents           []ExtractedIntent  `json:"intents" validate:"required"`
	Entities          []ExtractedEntity  `json:"entities"`
	SentimentAnalysis *SentimentAnalysis `json:"sentimentAnalysis"`
}

// Client is a basic LUIS client
type Client struct {
	httpClient  *http.Client
	httpRetries *httpx.RetryConfig
	httpAccess  *httpx.AccessConfig
	endpoint    string
	appID       string
	key         string
}

// NewClient creates a new client
func NewClient(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, endpoint, appID, key string) *Client {
	return &Client{
		httpClient:  httpClient,
		httpRetries: httpRetries,
		httpAccess:  httpAccess,
		endpoint:    endpoint,
		appID:       appID,
		key:         key,
	}
}

// Predict gets the published endpoint predictions for the given query
func (c *Client) Predict(q string) (*PredictResponse, *httpx.Trace, error) {
	endpoint := fmt.Sprintf("%s/apps/%s?verbose=true&subscription-key=%s&q=%s", c.endpoint, c.appID, c.key, url.QueryEscape(q))

	request, err := httpx.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	trace, err := httpx.DoTrace(c.httpClient, request, c.httpRetries, c.httpAccess, -1)
	if err != nil {
		return nil, trace, err
	}

	if trace.Response != nil && trace.Response.StatusCode == 200 {
		response := &PredictResponse{}
		if err := utils.UnmarshalAndValidate(trace.ResponseBody, response); err != nil {
			return nil, trace, err
		}
		return response, trace, nil
	}

	return nil, trace, errors.New("LUIS API request failed")
}
