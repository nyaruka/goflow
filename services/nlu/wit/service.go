package wit

import (
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// an NLU service implmentation for a wit.ai app
type service struct {
	classifier  *flows.Classifier
	accessToken string
}

// NewService creates a new NLU service
func NewService(classifier *flows.Classifier, accessToken string) flows.NLUService {
	return &service{
		classifier:  classifier,
		accessToken: accessToken,
	}
}

func (s *service) Classify(session flows.Session, input string, logEvent flows.EventCallback) (*flows.NLUClassification, error) {
	client := NewClient(session.Engine().HTTPClient(), s.accessToken)

	message, trace, err := client.Message(input)
	if trace != nil {
		status := flows.CallStatusSuccess
		if err != nil {
			status = flows.CallStatusResponseError
		}
		logEvent(events.NewClassifierCalled(
			s.classifier.Reference(),
			trace.Request.URL.String(),
			status,
			string(trace.RequestTrace),
			string(trace.ResponseTrace),
			int(trace.TimeTaken/time.Millisecond),
		))
	}
	if err != nil {
		return nil, err
	}

	result := &flows.NLUClassification{
		Intents:  make([]flows.ExtractedIntent, 0, 1),
		Entities: make(map[string][]flows.ExtractedEntity),
	}

	// wit returns intent as just another entity so we need to extract it by name
	for name, entity := range message.Entities {
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

var _ flows.NLUService = (*service)(nil)
