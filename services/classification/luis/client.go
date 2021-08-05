package luis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

type Intent struct {
	Score decimal.Decimal `json:"score"`
}

type Entity struct {
	Type   string          `json:"type"`
	Text   string          `json:"text" validate:"required"`
	Length int             `json:"length"`
	Score  decimal.Decimal `json:"score" validate:"required"`
}

type Entities struct {
	Values   map[string][]string
	Instance map[string][]Entity
}

func (e *Entities) UnmarshalJSON(data []byte) error {
	var asMap map[string]json.RawMessage
	if err := jsonx.Unmarshal(data, &asMap); err != nil {
		return err
	}

	// if $instance exists, parse it separately as it has a different structure
	if instanceJSON, hasInstance := asMap["$instance"]; hasInstance {
		if err := jsonx.Unmarshal(instanceJSON, &e.Instance); err != nil {
			return err
		}
		delete(asMap, "$instance")
	}

	// now we can parse the rest as string lists
	e.Values = make(map[string][]string, len(asMap))
	for key, data := range asMap {
		var v []string
		if err := jsonx.Unmarshal(data, &v); err != nil {
			return err
		}
		e.Values[key] = v
	}

	return nil
}

type Sentiment struct {
	Label string          `json:"label"`
	Score decimal.Decimal `json:"score"`
}

type Prediction struct {
	TopIntent string            `json:"topIntent"`
	Intents   map[string]Intent `json:"intents" validate:"required"`
	Entities  *Entities         `json:"entities"`
	Sentiment *Sentiment        `json:"sentiment"`
}

// PredictResponse is the response from a predict request
type PredictResponse struct {
	Query      string      `json:"query"`
	Prediction *Prediction `json:"prediction" validate:"required"`
}

// Client is a basic LUIS client
type Client struct {
	httpClient  *http.Client
	httpRetries *httpx.RetryConfig
	httpAccess  *httpx.AccessConfig
	endpoint    string
	appID       string
	key         string
	slot        string
}

// NewClient creates a new client
func NewClient(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, endpoint, appID, key, slot string) *Client {
	return &Client{
		httpClient:  httpClient,
		httpRetries: httpRetries,
		httpAccess:  httpAccess,
		endpoint:    endpoint,
		appID:       appID,
		key:         key,
		slot:        slot,
	}
}

// Predict gets the published endpoint predictions for the given query
func (c *Client) Predict(q string) (*PredictResponse, *httpx.Trace, error) {
	endpoint := fmt.Sprintf("%sluis/prediction/v3.0/apps/%s/slots/%s/predict?subscription-key=%s&verbose=true&show-all-intents=true&log=true&query=%s", c.endpoint, c.appID, c.slot, c.key, url.QueryEscape(q))

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
