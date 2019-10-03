package test

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/providers/webhooks"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// NewEngine creates an engine instance for testing
func NewEngine() flows.Engine {
	return engine.NewBuilder().
		WithWebhookService(webhooks.NewService("goflow-testing", 10000)).
		WithNLUService(func(s flows.Session, c assets.Classifier) flows.NLUProvider { return newNLUProvider(c) }).
		WithAirtimeService(func(flows.Session) flows.AirtimeProvider { return newAirtimeProvider("RWF") }).
		Build()
}

// implementation of NLUProvider for testing which always returns the first intent
type nluProvider struct {
	classifier assets.Classifier
}

func newNLUProvider(classifier assets.Classifier) *nluProvider {
	return &nluProvider{classifier: classifier}
}

func (p *nluProvider) Classify(session flows.Session, input string) (*flows.NLUClassification, error) {
	classifierIntents := p.classifier.Intents()
	extractedIntents := make([]flows.ExtractedIntent, len(p.classifier.Intents()))
	confidence := decimal.RequireFromString("0.5")
	for i := range classifierIntents {
		extractedIntents[i] = flows.ExtractedIntent{classifierIntents[i], confidence}
		confidence = confidence.Div(decimal.RequireFromString("2"))
	}

	return &flows.NLUClassification{
		Intents: extractedIntents,
		Entities: map[string][]flows.ExtractedEntity{
			"location": []flows.ExtractedEntity{
				flows.ExtractedEntity{"Quito", decimal.RequireFromString("1.0")},
			},
		},
	}, nil
}

var _ flows.NLUProvider = (*nluProvider)(nil)

// implementation of AirtimeProvider for testing which uses a fixed currency
type airtimeProvider struct {
	fixedCurrency string
}

func newAirtimeProvider(currency string) *airtimeProvider {
	return &airtimeProvider{fixedCurrency: currency}
}

func (p *airtimeProvider) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal) (*flows.AirtimeTransfer, error) {
	t := &flows.AirtimeTransfer{
		Sender:    sender,
		Recipient: recipient,
		Currency:  p.fixedCurrency,
		Status:    flows.AirtimeTransferStatusFailed,
	}

	amount, hasAmount := amounts[p.fixedCurrency]
	if !hasAmount {
		return t, errors.Errorf("no amount configured for transfers in %s", p.fixedCurrency)
	}

	t.DesiredAmount = amount
	t.ActualAmount = amount
	t.Status = flows.AirtimeTransferStatusSuccess
	return t, nil
}

var _ flows.AirtimeProvider = (*airtimeProvider)(nil)
