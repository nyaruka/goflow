package smtp

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/smtpx"

	"github.com/pkg/errors"
)

type service struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewService creates a new SMTP email service
func NewService(host string, port int, username, password, from string) flows.EmailService {
	return &service{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// NewServiceFromURL creates a new SMTP email service from a URL like smtp://user:pass@host:port/?from=from@example.com
func NewServiceFromURL(connectionURL string) (flows.EmailService, error) {
	url, err := url.Parse(connectionURL)
	if err != nil {
		return nil, errors.New("malformed connection URL")
	}
	if url.Scheme != "smtp" {
		return nil, errors.New("connection URL must use SMTP scheme")
	}

	host := url.Hostname()

	// parse port if provided or default to 25
	port := 25
	if url.Port() != "" {
		port, err = strconv.Atoi(url.Port())
		if err != nil || port < 0 || port > 65535 {
			return nil, errors.Errorf("%s is not a valid port number", url.Port())
		}
	}

	// get the credentials
	if url.User == nil {
		return nil, errors.New("missing credentials in connection URL")
	}
	username := url.User.Username()
	password, _ := url.User.Password()

	// get our from address
	from := url.Query().Get("from")
	if from == "" {
		from = fmt.Sprintf("%s@%s", username, host) // default to username@host if not set
	}

	return NewService(host, port, username, password, from), nil
}

func (s *service) Send(session flows.Session, addresses []string, subject, body string) error {
	return smtpx.Send(s.host, s.port, s.username, s.password, s.from, addresses, subject, body)
}
