package zendesk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/jsonx"
)

// Client is a basic zendesk client
type Client struct {
	httpClient  *http.Client
	httpRetries *httpx.RetryConfig
	subdomain   string
	username    string
	apiToken    string
}

// NewClient creates a new zendesk client
func NewClient(httpClient *http.Client, httpRetries *httpx.RetryConfig, subdomain, username, apiToken string) *Client {
	return &Client{
		httpClient:  httpClient,
		httpRetries: httpRetries,
		subdomain:   subdomain,
		username:    username,
		apiToken:    apiToken,
	}
}

type errorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

type ticketComment struct {
	Body string `json:"body"`
}

type newTicket struct {
	Subject    string        `json:"subject"`
	Comment    ticketComment `json:"comment"`
	ExternalID string        `json:"external_id"`
}

// Ticket is a ticket in Zendesk
type Ticket struct {
	ID         int       `json:"id"`
	URL        string    `json:"url"`
	ExternalID string    `json:"external_id"`
	CreatedAt  time.Time `json:"created_at"`
	Subject    string    `json:"subject"`
}

// CreateTicket creates a new ticket
func (c *Client) CreateTicket(subject, body string) (*Ticket, *httpx.Trace, error) {
	r := struct {
		Ticket newTicket `json:"ticket"`
	}{
		Ticket: newTicket{
			Subject: subject,
			Comment: ticketComment{Body: body},
		},
	}

	trace, err := c.post("tickets", &r)
	if err != nil {
		return nil, trace, err
	}

	if trace.Response.StatusCode >= 400 {
		response := &errorResponse{}
		jsonx.Unmarshal(trace.ResponseBody, response)
		return nil, trace, errors.New(response.Description)
	}

	response := struct {
		Ticket Ticket `json:"ticket"`
	}{}
	if err := jsonx.Unmarshal(trace.ResponseBody, &response); err != nil {
		return nil, trace, err
	}

	return &response.Ticket, trace, nil
}

func (c *Client) post(endpoint string, payload interface{}) (*httpx.Trace, error) {
	data, err := jsonx.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s.zendesk.com/api/v2/%s.json", c.subdomain, endpoint), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.username, c.apiToken)

	return httpx.DoTrace(c.httpClient, req, c.httpRetries, nil, -1)
}
