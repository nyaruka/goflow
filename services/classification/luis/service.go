package luis

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils/httpx"
)

// a classification service implmentation for a LUIS app
type service struct {
	classifier *flows.Classifier
	endpoint   string
	appID      string
	key        string
}

// NewService creates a new classification service
func NewService(classifier *flows.Classifier, endpoint, appID, key string) flows.ClassificationService {
	return &service{
		classifier: classifier,
		endpoint:   endpoint,
		appID:      appID,
		key:        key,
	}
}

func (s *service) Classify(session flows.Session, input string, logEvent flows.EventCallback) (*flows.Classification, error) {
	client := NewClient(session.Engine().HTTPClient(), s.endpoint, s.appID, s.key)

	response, trace, err := client.Predict(input)
	if trace != nil {
		logEvent(events.NewClassifierCalled(s.classifier.Reference(), []*httpx.Trace{trace}))
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
