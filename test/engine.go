package test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/shopspring/decimal"
)

// NewEngine creates an engine instance for testing
func NewEngine() flows.Engine {
	retries := httpx.NewFixedRetries(1*time.Millisecond, 2*time.Millisecond)

	return engine.NewBuilder().
		WithMaxFieldChars(256).
		WithEmailServiceFactory(func(s flows.SessionAssets) (flows.EmailService, error) {
			return newEmailService(), nil
		}).
		WithWebhookServiceFactory(webhooks.NewServiceFactory(http.DefaultClient, retries, nil, map[string]string{"User-Agent": "goflow-testing"}, 10000)).
		WithClassificationServiceFactory(func(c *flows.Classifier) (flows.ClassificationService, error) {
			return newClassificationService(c), nil
		}).
		WithLLMServiceFactory(func(l *flows.LLM) (flows.LLMService, error) {
			return newLLMService(l), nil
		}).
		WithAirtimeServiceFactory(func(flows.SessionAssets) (flows.AirtimeService, error) { return newAirtimeService("RWF"), nil }).
		Build()
}

// implementation of an email service for testing which just fakes sending the email
type emailService struct{}

func newEmailService() *emailService {
	return &emailService{}
}

func (s *emailService) Send(addresses []string, subject, body string) error {
	return nil
}

// implementation of a classification service for testing which always returns the first intent
type classificationService struct {
	classifier *flows.Classifier
}

func newClassificationService(classifier *flows.Classifier) *classificationService {
	return &classificationService{classifier: classifier}
}

func (s *classificationService) Classify(env envs.Environment, input string, logHTTP flows.HTTPLogCallback) (*flows.Classification, error) {
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

var _ flows.ClassificationService = (*classificationService)(nil)

// implementation of an LLM service for testing which returns the last word of the input
type llmService struct {
	llm *flows.LLM
}

func newLLMService(llm *flows.LLM) *llmService {
	return &llmService{llm: llm}
}

func (s *llmService) Response(ctx context.Context, env envs.Environment, instructions, input string) (string, error) {
	// get last word from the instructions and return that.. because we're not a real LLM!
	words := strings.Fields(instructions)
	output := words[len(words)-1]

	return output, nil
}

var _ flows.LLMService = (*llmService)(nil)

// implementation of an airtime service for testing which uses a fixed currency
type airtimeService struct {
	fixedCurrency string
}

func newAirtimeService(currency string) *airtimeService {
	return &airtimeService{fixedCurrency: currency}
}

func (s *airtimeService) Transfer(sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logHTTP flows.HTTPLogCallback) (*flows.AirtimeTransfer, error) {
	logHTTP(&flows.HTTPLog{
		HTTPLogWithoutTime: &flows.HTTPLogWithoutTime{
			LogWithoutTime: &httpx.LogWithoutTime{
				URL:        "http://send.airtime.com",
				StatusCode: 200,
				Request:    "GET / HTTP/1.1\r\nHost: send.airtime.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
				Response:   "HTTP/1.0 200 OK\r\nContent-Length: 15\r\n\r\n{\"status\":\"ok\"}",
				ElapsedMS:  0,
				Retries:    0,
			},
			Status: flows.CallStatusSuccess,
		},
		CreatedOn: time.Date(2019, 10, 16, 13, 59, 30, 123456789, time.UTC),
	})

	amount, hasAmount := amounts[s.fixedCurrency]
	if !hasAmount {
		return nil, fmt.Errorf("no amount configured for transfers in %s", s.fixedCurrency)
	}

	transfer := &flows.AirtimeTransfer{
		UUID:      flows.AirtimeTransferUUID(uuids.NewV4()),
		Sender:    sender,
		Recipient: recipient,
		Currency:  s.fixedCurrency,
		Amount:    amount,
	}

	return transfer, nil
}

var _ flows.AirtimeService = (*airtimeService)(nil)
