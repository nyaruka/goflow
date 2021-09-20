package test

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/services/webhooks"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// NewEngine creates an engine instance for testing
func NewEngine() flows.Engine {
	retries := httpx.NewFixedRetries(1*time.Millisecond, 2*time.Millisecond)

	return engine.NewBuilder().
		WithEmailServiceFactory(func(s flows.Session) (flows.EmailService, error) {
			return newEmailService(), nil
		}).
		WithWebhookServiceFactory(webhooks.NewServiceFactory(http.DefaultClient, retries, nil, map[string]string{"User-Agent": "goflow-testing"}, 10000)).
		WithClassificationServiceFactory(func(s flows.Session, c *flows.Classifier) (flows.ClassificationService, error) {
			return newClassificationService(c), nil
		}).
		WithTicketServiceFactory(func(s flows.Session, t *flows.Ticketer) (flows.TicketService, error) { return NewTicketService(t), nil }).
		WithAirtimeServiceFactory(func(flows.Session) (flows.AirtimeService, error) { return newAirtimeService("RWF"), nil }).
		Build()
}

// implementation of an email service for testing which just fakes sending the email
type emailService struct{}

func newEmailService() *emailService {
	return &emailService{}
}

func (s *emailService) Send(session flows.Session, addresses []string, subject, body string) error {
	return nil
}

// implementation of a classification service for testing which always returns the first intent
type classificationService struct {
	classifier *flows.Classifier
}

func newClassificationService(classifier *flows.Classifier) *classificationService {
	return &classificationService{classifier: classifier}
}

func (s *classificationService) Classify(session flows.Session, input string, logHTTP flows.HTTPLogCallback) (*flows.Classification, error) {
	classifierIntents := s.classifier.Intents()
	extractedIntents := make([]flows.ExtractedIntent, len(s.classifier.Intents()))
	confidence := decimal.RequireFromString("0.5")
	for i := range classifierIntents {
		extractedIntents[i] = flows.ExtractedIntent{Name: classifierIntents[i], Confidence: confidence}
		confidence = confidence.Div(decimal.RequireFromString("2"))
	}

	logHTTP(&flows.HTTPLog{
		HTTPTrace: &flows.HTTPTrace{
			URL:        "http://test.acme.ai?classify",
			StatusCode: 200,
			Status:     flows.CallStatusSuccess,
			Request:    "GET /?classify HTTP/1.1\r\nHost: test.acme.ai\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
			Response:   "HTTP/1.0 200 OK\r\nContent-Length: 14\r\n\r\n{\"intents\":[]}",
			ElapsedMS:  1000,
			Retries:    0,
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

// implementation of a ticket service for testing which fails if ticket subject contains "fail" and passes if not
type ticketService struct {
	ticketer *flows.Ticketer
}

// NewTicketService creates a new ticket service for testing
func NewTicketService(ticketer *flows.Ticketer) flows.TicketService {
	return &ticketService{ticketer: ticketer}
}

func (s *ticketService) Open(session flows.Session, topic *flows.Topic, body string, assignee *flows.User, logHTTP flows.HTTPLogCallback) (*flows.Ticket, error) {
	if strings.Contains(body, "fail") {
		logHTTP(&flows.HTTPLog{
			HTTPTrace: &flows.HTTPTrace{
				URL:        "http://nyaruka.tickets.com/tickets.json",
				StatusCode: 400,
				Status:     flows.CallStatusResponseError,
				Request:    fmt.Sprintf("POST /tickets.json HTTP/1.1\r\nAccept-Encoding: gzip\r\n\r\n{\"body\":\"%s\"}", body),
				Response:   "HTTP/1.0 400 OK\r\nContent-Length: 17\r\n\r\n{\"status\":\"fail\"}",
				ElapsedMS:  1,
				Retries:    0,
			},
			CreatedOn: time.Date(2019, 10, 16, 13, 59, 30, 123456789, time.UTC),
		})

		return nil, errors.New("error calling ticket API")
	}

	logHTTP(&flows.HTTPLog{
		HTTPTrace: &flows.HTTPTrace{
			URL:        "http://nyaruka.tickets.com/tickets.json",
			StatusCode: 200,
			Status:     flows.CallStatusSuccess,
			Request:    fmt.Sprintf("POST /tickets.json HTTP/1.1\r\nAccept-Encoding: gzip\r\n\r\n{\"body\":\"%s\"}", body),
			Response:   "HTTP/1.0 200 OK\r\nContent-Length: 15\r\n\r\n{\"status\":\"ok\"}",
			ElapsedMS:  1,
			Retries:    0,
		},
		CreatedOn: time.Date(2019, 10, 16, 13, 59, 30, 123456789, time.UTC),
	})

	ticket := flows.OpenTicket(s.ticketer, topic, body, assignee)
	ticket.SetExternalID("123456")
	return ticket, nil
}

// implementation of an airtime service for testing which uses a fixed currency
type airtimeService struct {
	fixedCurrency string
}

func newAirtimeService(currency string) *airtimeService {
	return &airtimeService{fixedCurrency: currency}
}

func (s *airtimeService) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logHTTP flows.HTTPLogCallback) (*flows.AirtimeTransfer, error) {
	logHTTP(&flows.HTTPLog{
		HTTPTrace: &flows.HTTPTrace{
			URL:        "http://send.airtime.com",
			StatusCode: 200,
			Status:     flows.CallStatusSuccess,
			Request:    "GET / HTTP/1.1\r\nHost: send.airtime.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
			Response:   "HTTP/1.0 200 OK\r\nContent-Length: 15\r\n\r\n{\"status\":\"ok\"}",
			ElapsedMS:  0,
			Retries:    0,
		},
		CreatedOn: time.Date(2019, 10, 16, 13, 59, 30, 123456789, time.UTC),
	})

	amount, hasAmount := amounts[s.fixedCurrency]
	if !hasAmount {
		return nil, errors.Errorf("no amount configured for transfers in %s", s.fixedCurrency)
	}

	transfer := &flows.AirtimeTransfer{
		Sender:        sender,
		Recipient:     recipient,
		Currency:      s.fixedCurrency,
		DesiredAmount: amount,
		ActualAmount:  amount,
	}

	return transfer, nil
}

var _ flows.AirtimeService = (*airtimeService)(nil)
