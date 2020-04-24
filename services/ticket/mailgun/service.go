package mailgun

import (
	"fmt"
	"net/http"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/httpx"
)

type service struct {
	client   *Client
	ticketer *flows.Ticketer
	to       string
}

// NewService creates a new Mailgun email-based ticketing service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, ticketer *flows.Ticketer, domain, token, to string) flows.TicketService {
	return &service{
		client:   NewClient(httpClient, httpRetries, domain, token),
		ticketer: ticketer,
		to:       to,
	}
}

func (s *service) Open(session flows.Session, subject, body string, logHTTP flows.HTTPLogCallback) (*flows.Ticket, error) {
	ticket := flows.NewTicket(s.ticketer, subject, body)

	fromAddress := fmt.Sprintf("thread+%s@%s", ticket.UUID, s.client.domain)

	trace, err := s.client.SendMessage(fromAddress, s.to, subject, body)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode))
	}
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
