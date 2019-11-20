package luis

import (
	"net/http"

	"github.com/nyaruka/goflow/flows"
)

// a classification service implementation for a LUIS app
type service struct {
	httpClient *http.Client
	classifier *flows.Classifier
	endpoint   string
	appID      string
	key        string
}

// NewService creates a new classification service
func NewService(httpClient *http.Client, classifier *flows.Classifier, endpoint, appID, key string) flows.ClassificationService {
	return &service{
		httpClient: httpClient,
		classifier: classifier,
		endpoint:   endpoint,
		appID:      appID,
		key:        key,
	}
}

func (s *service) Classify(session flows.Session, input string, logHTTP flows.HTTPLogCallback) (*flows.Classification, error) {
	client := NewClient(s.httpClient, s.endpoint, s.appID, s.key)

	response, trace, err := client.Predict(input)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode))
	}
	if err != nil {
		return nil, err
	}

	result := &flows.Classification{
		Intents:  make([]flows.ExtractedIntent, len(response.Intents)),
		Entities: make(map[string][]flows.ExtractedEntity, len(response.Entities)),
	}

	for i, intent := range response.Intents {
		result.Intents[i] = flows.ExtractedIntent{Name: intent.Intent, Confidence: intent.Score}
	}

	for _, entity := range response.Entities {
		result.Entities[entity.Type] = []flows.ExtractedEntity{
			flows.ExtractedEntity{Value: entity.Entity, Confidence: entity.Score},
		}
	}

	// if sentiment analysis was included, convert to an entity
	if response.SentimentAnalysis != nil {
		result.Entities["sentiment"] = []flows.ExtractedEntity{
			flows.ExtractedEntity{Value: response.SentimentAnalysis.Label, Confidence: response.SentimentAnalysis.Score},
		}
	}

	return result, nil
}

var _ flows.ClassificationService = (*service)(nil)
