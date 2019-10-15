package test

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// NewEngine creates an engine instance for testing
func NewEngine() flows.Engine {
	return engine.NewBuilder().
		WithWebhookServiceFactory(webhooks.NewServiceFactory("goflow-testing", 10000)).
		WithClassificationServiceFactory(func(s flows.Session, c *flows.Classifier) (flows.ClassificationService, error) {
			return newClassificationService(c), nil
		}).
		WithAirtimeServiceFactory(func(flows.Session) (flows.AirtimeService, error) { return NewAirtimeService("RWF"), nil }).
		Build()
}

// implementation of an NLU service for testing which always returns the first intent
type nluService struct {
	classifier *flows.Classifier
}

func newClassificationService(classifier *flows.Classifier) *nluService {
	return &nluService{classifier: classifier}
}

func (s *nluService) Classify(session flows.Session, input string, logEvent flows.EventCallback) (*flows.Classification, error) {
	classifierIntents := s.classifier.Intents()
	extractedIntents := make([]flows.ExtractedIntent, len(s.classifier.Intents()))
	confidence := decimal.RequireFromString("0.5")
	for i := range classifierIntents {
		extractedIntents[i] = flows.ExtractedIntent{classifierIntents[i], confidence}
		confidence = confidence.Div(decimal.RequireFromString("2"))
	}

	logEvent(events.NewClassifierCalled(
		s.classifier.Reference(),
		[]*httpx.Trace{
			httpx.NewMockTrace("GET", "http://test.acme.ai?classifiy", 200, `{"intents":[]}`),
		},
	))

	classification := &flows.Classification{
		Intents: extractedIntents,
		Entities: map[string][]flows.ExtractedEntity{
			"location": []flows.ExtractedEntity{
				flows.ExtractedEntity{"Quito", decimal.RequireFromString("1.0")},
			},
		},
	}

	return classification, nil
}

var _ flows.ClassificationService = (*nluService)(nil)

// implementation of an airtime service for testing which uses a fixed currency
type airtimeService struct {
	fixedCurrency string
}

func NewAirtimeService(currency string) flows.AirtimeService {
	return &airtimeService{fixedCurrency: currency}
}

func (s *airtimeService) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logEvent flows.EventCallback) (*flows.AirtimeTransfer, error) {
	amount, hasAmount := amounts[s.fixedCurrency]
	if !hasAmount {
		return nil, errors.Errorf("no amount configured for transfers in %s", s.fixedCurrency)
	}

	transfer := &flows.AirtimeTransfer{
		Sender:    sender,
		Recipient: recipient,
		Currency:  s.fixedCurrency,
		Amount:    amount,
	}

	logEvent(events.NewAirtimeTransferred(
		transfer,
		[]*httpx.Trace{
			httpx.NewMockTrace("GET", "http://send.airtime.com", 200, `{"status":"ok"}`),
		},
	))

	return transfer, nil
}

var _ flows.AirtimeService = (*airtimeService)(nil)
