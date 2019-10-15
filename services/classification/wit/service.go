package wit

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils/httpx"
)

// a classification service implmentation for a wit.ai app
type service struct {
	classifier  *flows.Classifier
	accessToken string
}

// NewService creates a new classification service
func NewService(classifier *flows.Classifier, accessToken string) flows.ClassificationService {
	return &service{
		classifier:  classifier,
		accessToken: accessToken,
	}
}

func (s *service) Classify(session flows.Session, input string, logEvent flows.EventCallback) (*flows.Classification, error) {
	client := NewClient(session.Engine().HTTPClient(), s.accessToken)

	response, trace, err := client.Message(input)
	if trace != nil {
		logEvent(events.NewClassifierCalled(s.classifier.Reference(), []*httpx.Trace{trace}))
	}
	if err != nil {
		return nil, err
	}

	result := &flows.Classification{
		Intents:  make([]flows.ExtractedIntent, 0, 1),
		Entities: make(map[string][]flows.ExtractedEntity),
	}

	// wit returns intent as just another entity so we need to extract it by name
	for name, entity := range response.Entities {
		if name == "intent" {
			for _, candidate := range entity {
				result.Intents = append(result.Intents, flows.ExtractedIntent{
					Name:       candidate.Value,
					Confidence: candidate.Confidence,
				})
			}
		} else {
			entities := make([]flows.ExtractedEntity, 0, len(entity))
			for _, candidate := range entity {
				entities = append(entities, flows.ExtractedEntity{
					Value:      candidate.Value,
					Confidence: candidate.Confidence,
				})
			}
			result.Entities[name] = entities
		}
	}

	return result, nil
}

var _ flows.ClassificationService = (*service)(nil)
