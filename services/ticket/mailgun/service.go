package mailgun

import (
	"fmt"
	"net/http"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/pkg/errors"
)

type service struct {
	client    *Client
	ticketer  *flows.Ticketer
	toAddress string
}

// NewService creates a new Mailgun email-based ticketing service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, ticketer *flows.Ticketer, domain, apiKey, toAddress string) flows.TicketService {
	return &service{
		client:    NewClient(httpClient, httpRetries, domain, apiKey),
		ticketer:  ticketer,
		toAddress: toAddress,
	}
}

// Open opens a ticket which for mailgun means just sending an initial email
func (s *service) Open(session flows.Session, subject, body string, logHTTP flows.HTTPLogCallback) (*flows.Ticket, error) {
	ticketUUID := flows.TicketUUID(uuids.New())

	fromAddress := fmt.Sprintf("ticket+%s@%s", ticketUUID, s.client.domain)
	from := fmt.Sprintf("%s <%s>", session.Contact().Format(session.Environment()), fromAddress)

	_, trace, err := s.client.SendMessage(from, s.toAddress, subject, body)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode))
	}
	if err != nil {
		return nil, errors.Wrap(err, "error calling mailgun API")
	}

	return flows.NewTicket(ticketUUID, s.ticketer, subject, body, fromAddress), nil
}
