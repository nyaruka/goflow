package luis

import (
	"net/http"
	"sort"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// a classification service implementation for a LUIS app
type service struct {
	client     *Client
	classifier *flows.Classifier
	redactor   utils.Redactor
}

// NewService creates a new classification service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, classifier *flows.Classifier, endpoint, appID, key, slot string) flows.ClassificationService {
	return &service{
		client:     NewClient(httpClient, httpRetries, httpAccess, endpoint, appID, key, slot),
		classifier: classifier,
		redactor:   utils.NewRedactor(flows.RedactionMask, key),
	}
}

func (s *service) Classify(session flows.Session, input string, logHTTP flows.HTTPLogCallback) (*flows.Classification, error) {
	response, trace, err := s.client.Predict(input)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode, s.redactor))
	}
	if err != nil {
		return nil, err
	}

	result := &flows.Classification{
		Intents:  make([]flows.ExtractedIntent, 0, len(response.Prediction.Intents)),
		Entities: make(map[string][]flows.ExtractedEntity, len(response.Prediction.Entities.Values)),
	}

	for name, intent := range response.Prediction.Intents {
		result.Intents = append(result.Intents, flows.ExtractedIntent{Name: name, Confidence: intent.Score})
	}
	sort.SliceStable(result.Intents, func(i, j int) bool { return result.Intents[i].Confidence.GreaterThan(result.Intents[j].Confidence) })

	for name, matches := range response.Prediction.Entities.Instance {
		var entities []flows.ExtractedEntity
		for _, match := range matches {
			entities = append(entities, flows.ExtractedEntity{Value: match.Text, Confidence: match.Score})
		}
		result.Entities[name] = entities
	}

	// if sentiment analysis was included, convert to an entity
	if response.Prediction.Sentiment != nil {
		result.Entities["sentiment"] = []flows.ExtractedEntity{
			{Value: response.Prediction.Sentiment.Label, Confidence: response.Prediction.Sentiment.Score},
		}
	}

	return result, nil
}

var _ flows.ClassificationService = (*service)(nil)
