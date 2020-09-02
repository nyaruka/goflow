package luis

import (
	"net/http"

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
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, classifier *flows.Classifier, endpoint, appID, key string) flows.ClassificationService {
	return &service{
		client:     NewClient(httpClient, httpRetries, httpAccess, endpoint, appID, key),
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
		Intents:  make([]flows.ExtractedIntent, len(response.Intents)),
		Entities: make(map[string][]flows.ExtractedEntity, len(response.Entities)),
	}

	for i, intent := range response.Intents {
		result.Intents[i] = flows.ExtractedIntent{Name: intent.Intent, Confidence: intent.Score}
	}

	for _, entity := range response.Entities {
		result.Entities[entity.Type] = []flows.ExtractedEntity{
			{Value: entity.Entity, Confidence: entity.Score},
		}
	}

	// if sentiment analysis was included, convert to an entity
	if response.SentimentAnalysis != nil {
		result.Entities["sentiment"] = []flows.ExtractedEntity{
			{Value: response.SentimentAnalysis.Label, Confidence: response.SentimentAnalysis.Score},
		}
	}

	return result, nil
}

var _ flows.ClassificationService = (*service)(nil)
