package bothub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/httpx"

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
	Intent        IntentMatch                         `json:"intent" validate:"required"`
	IntentRanking []IntentMatch                       `json:"intent_ranking" validate:"required"`
	LabelsList    []string                            `json:"labels_list"`
	EntitiesList  []string                            `json:"entities_list"`
	Entities      map[string]map[string][]EntityMatch `json:"entities"`
	Text          string                              `json:"text"`
	UpdateID      int                                 `json:"update_id"`
	Language      string                              `json:"language"`
}

// Client is a basic Wit.ai client
type Client struct {
	httpClient  *http.Client
	accessToken string
}

// NewClient creates a new client
func NewClient(httpClient *http.Client, accessToken string) *Client {
	return &Client{
		httpClient:  httpClient,
		accessToken: accessToken,
	}
}

// Parse does a parse of the given text
func (c *Client) Parse(text string) (*ParseResponse, *httpx.Trace, error) {
	endpoint := fmt.Sprintf("%s/parse", apiBaseURL)

	form := url.Values{}
	form.Add("text", text)

	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": fmt.Sprintf("Bearer %s", c.accessToken),
	}

	trace, err := httpx.DoTrace(c.httpClient, "POST", endpoint, strings.NewReader(form.Encode()), headers)
	if err != nil {
		return nil, nil, err
	}

	response := &ParseResponse{}
	if err := utils.UnmarshalAndValidate(trace.Body, response); err != nil {
		return nil, trace, err
	}

	return response, trace, nil
}
