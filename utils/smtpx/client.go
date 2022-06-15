package smtpx

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/Shopify/gomail"
	"github.com/pkg/errors"
)

// Client is an SMTP client
type Client struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewClient creates a new client
func NewClient(host string, port int, username, password, from string) *Client {
	return &Client{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// NewClientFromURL creates a client from a URL like smtp://user:pass@host:port/?from=from@example.com
func NewClientFromURL(connectionURL string) (*Client, error) {
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

	return NewClient(host, port, username, password, from), nil
}

// Send sends the given message
func (c *Client) Send(m *Message) error {
	// create MIME message
	mm := gomail.NewMessage()
	mm.SetHeader("From", c.from)
	mm.SetHeader("To", m.recipients...)
	mm.SetHeader("Subject", m.subject)
	mm.SetBody("text/plain", m.text)
	if m.html != "" {
		mm.AddAlternative("text/html", m.html)
	}

	d := gomail.NewDialer(c.host, c.port, c.username, c.password)
	return d.DialAndSend(mm)
}
