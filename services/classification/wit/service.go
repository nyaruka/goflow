package wit

import (
	"net/http"
	"strings"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// a classification service implementation for a wit.ai app
type service struct {
	client     *Client
	classifier *flows.Classifier
	redactor   utils.Redactor
}

// NewService creates a new classification service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, classifier *flows.Classifier, accessToken string) flows.ClassificationService {
	return &service{
		client:     NewClient(httpClient, httpRetries, accessToken),
		classifier: classifier,
		redactor:   utils.NewRedactor(flows.RedactionMask, accessToken),
	}
}

func (s *service) Classify(session flows.Session, input string, logHTTP flows.HTTPLogCallback) (*flows.Classification, error) {
	response, trace, err := s.client.Message(input)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode, s.redactor))
	}
	if err != nil {
		return nil, err
	}

	result := &flows.Classification{
		Intents:  make([]flows.ExtractedIntent, len(response.Intents)),
		Entities: make(map[string][]flows.ExtractedEntity),
	}

	for i, intent := range response.Intents {
		result.Intents[i] = flows.ExtractedIntent{Name: intent.Name, Confidence: intent.Confidence}
	}

	for nameAndRole, entity := range response.Entities {
		name := strings.Split(nameAndRole, ":")[0]
		entities := make([]flows.ExtractedEntity, 0, len(entity))
		for _, candidate := range entity {
			entities = append(entities, flows.ExtractedEntity{
				Value:      candidate.Value,
				Confidence: candidate.Confidence,
			})
		}
		result.Entities[name] = entities
	}

	return result, nil
}

var _ flows.ClassificationService = (*service)(nil)
