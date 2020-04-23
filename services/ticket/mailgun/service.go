package mailgun

import (
	"net/http"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"
)

type service struct {
	httpClient  *http.Client
	httpRetries *httpx.RetryConfig
	address     string
	token       string
}

// NewService creates a new Mailgun email-based ticketing service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, address, token string) flows.TicketService {
	return &service{
		httpClient:  httpClient,
		httpRetries: httpRetries,
		address:     address,
		token:       token,
	}
}

func (s *service) Open(session flows.Session, subject string, logHTTP flows.HTTPLogCallback) (*flows.Ticket, error) {
	return &flows.Ticket{ID: string(uuids.New()), Subject: subject}, nil
}
