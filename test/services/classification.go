package services

import (
	"time"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/shopspring/decimal"
)

// implementation of a classification service for testing which always returns the first intent
type Classification struct {
	classifier *flows.Classifier
}

func NewClassification(classifier *flows.Classifier) *Classification {
	return &Classification{classifier: classifier}
}

func (s *Classification) Classify(env envs.Environment, input string, logHTTP flows.HTTPLogCallback) (*flows.Classification, error) {
	classifierIntents := s.classifier.Intents()
	extractedIntents := make([]flows.ExtractedIntent, len(s.classifier.Intents()))
	confidence := decimal.RequireFromString("0.5")
	for i := range classifierIntents {
		extractedIntents[i] = flows.ExtractedIntent{Name: classifierIntents[i], Confidence: confidence}
		confidence = confidence.Div(decimal.RequireFromString("2"))
	}

	logHTTP(&flows.HTTPLog{
		HTTPLogWithoutTime: &flows.HTTPLogWithoutTime{
			LogWithoutTime: &httpx.LogWithoutTime{
				URL:        "http://test.acme.ai?classify",
				StatusCode: 200,
				Request:    "GET /?classify HTTP/1.1\r\nHost: test.acme.ai\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
				Response:   "HTTP/1.0 200 OK\r\nContent-Length: 14\r\n\r\n{\"intents\":[]}",
				ElapsedMS:  1000,
				Retries:    0,
			},
			Status: flows.CallStatusSuccess,
		},
		CreatedOn: time.Date(2019, 10, 16, 13, 59, 30, 123456789, time.UTC),
	})

	classification := &flows.Classification{
		Intents: extractedIntents,
		Entities: map[string][]flows.ExtractedEntity{
			"location": {
				flows.ExtractedEntity{Value: "Quito", Confidence: decimal.RequireFromString("1.0")},
			},
		},
	}

	return classification, nil
}

var _ flows.ClassificationService = (*Classification)(nil)
