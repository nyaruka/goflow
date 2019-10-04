package test

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/services/webhooks"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// NewEngine creates an engine instance for testing
func NewEngine() flows.Engine {
	return engine.NewBuilder().
		WithWebhookServiceFactory(webhooks.NewServiceFactory("goflow-testing", 10000)).
		WithNLUServiceFactory(func(s flows.Session, c *flows.Classifier) flows.NLUService { return newNLUService(c) }).
		WithAirtimeServiceFactory(func(flows.Session) flows.AirtimeService { return newAirtimeService("RWF") }).
		Build()
}

// implementation of an NLU service for testing which always returns the first intent
type nluService struct {
	classifier *flows.Classifier
}

func newNLUService(classifier *flows.Classifier) *nluService {
	return &nluService{classifier: classifier}
}

func (s *nluService) Classify(session flows.Session, input string, logEvent flows.EventCallback) (*flows.NLUClassification, error) {
	classifierIntents := s.classifier.Intents()
	extractedIntents := make([]flows.ExtractedIntent, len(s.classifier.Intents()))
	confidence := decimal.RequireFromString("0.5")
	for i := range classifierIntents {
		extractedIntents[i] = flows.ExtractedIntent{classifierIntents[i], confidence}
		confidence = confidence.Div(decimal.RequireFromString("2"))
	}

	logEvent(events.NewClassifierCalled(
		s.classifier.Reference(),
		"http://test.acme.ai?classifiy",
		flows.CallStatusSuccess,
		"GET /message?v=20170307&q=hello HTTP/1.1",
		"HTTP/1.1 200 OK\r\n\r\n{\"intents\":[]}",
		1,
	))

	classification := &flows.NLUClassification{
		Intents: extractedIntents,
		Entities: map[string][]flows.ExtractedEntity{
			"location": []flows.ExtractedEntity{
				flows.ExtractedEntity{"Quito", decimal.RequireFromString("1.0")},
			},
		},
	}

	return classification, nil
}

var _ flows.NLUService = (*nluService)(nil)

// implementation of an airtime service for testing which uses a fixed currency
type airtimeService struct {
	fixedCurrency string
}

func newAirtimeService(currency string) *airtimeService {
	return &airtimeService{fixedCurrency: currency}
}

func (s *airtimeService) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal) (*flows.AirtimeTransfer, error) {
	t := &flows.AirtimeTransfer{
		Sender:    sender,
		Recipient: recipient,
		Currency:  s.fixedCurrency,
		Status:    flows.AirtimeTransferStatusFailed,
	}

	amount, hasAmount := amounts[s.fixedCurrency]
	if !hasAmount {
		return t, errors.Errorf("no amount configured for transfers in %s", s.fixedCurrency)
	}

	t.DesiredAmount = amount
	t.ActualAmount = amount
	t.Status = flows.AirtimeTransferStatusSuccess
	return t, nil
}

var _ flows.AirtimeService = (*airtimeService)(nil)
