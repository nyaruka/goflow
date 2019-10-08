package luis

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/shopspring/decimal"
)

const (
	apiBaseURL = "https://westus.api.cognitive.microsoft.com/luis/v2.0"
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
	Query             string            `json:"query"`
	TopScoringIntent  *ExtractedIntent  `json:"topScoringIntent"`
	Intents           []ExtractedIntent `json:"intents" validate:"required"`
	Entities          []ExtractedEntity `json:"entities"`
	SentimentAnalysis SentimentAnalysis `json:"sentimentAnalysis"`
}

// Client is a basic LUIS client
type Client struct {
	httpClient *http.Client
	appID      string
	key        string
}

// NewClient creates a new client
func NewClient(httpClient *http.Client, appID, key string) *Client {
	return &Client{
		httpClient: httpClient,
		appID:      appID,
		key:        key,
	}
}

// Predict gets the published endpoint predictions for the given query
func (c *Client) Predict(q string) (*PredictResponse, *httpx.Trace, error) {
	endpoint := fmt.Sprintf("%s/apps/%s/?verbose=true&subscription-key=%s&q=%s", apiBaseURL, c.appID, c.key, url.QueryEscape(q))

	trace, err := httpx.DoTrace(c.httpClient, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	response := &PredictResponse{}
	if err := utils.UnmarshalAndValidate(trace.Body, response); err != nil {
		return nil, trace, err
	}

	return response, trace, nil
}
