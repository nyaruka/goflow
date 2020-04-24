package mailgun

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"
)

const apiBaseURL = "https://api.mailgun.net/v3"

// Client is a basic mailgun client
type Client struct {
	httpClient  *http.Client
	httpRetries *httpx.RetryConfig
	domain      string
	apiKey      string
}

// NewClient creates a new mailgun client
func NewClient(httpClient *http.Client, httpRetries *httpx.RetryConfig, domain, apiKey string) *Client {
	return &Client{
		httpClient:  httpClient,
		httpRetries: httpRetries,
		domain:      domain,
		apiKey:      apiKey,
	}
}

// SendMessage sends a new email message
func (c *Client) SendMessage(from, to, subject, text string) (*httpx.Trace, error) {
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	w.SetBoundary(string(uuids.New()))
	w.WriteField("from", from)
	w.WriteField("to", to)
	w.WriteField("subject", subject)
	w.WriteField("text", text)
	w.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/messages", apiBaseURL, c.domain), bytes.NewReader(b.Bytes()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("api", c.apiKey)

	return httpx.DoTrace(c.httpClient, req, c.httpRetries, nil, -1)
}
