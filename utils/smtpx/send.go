package smtpx

import (
	"time"
)

// Message is email message
type Message struct {
	recipients []string
	subject    string
	text       string
	html       string
}

// NewMessage creates a new message
func NewMessage(recipients []string, subject, text, html string) *Message {
	return &Message{
		recipients: recipients,
		subject:    subject,
		text:       text,
		html:       html,
	}
}

// Send an email using SMTP
func Send(c *Client, m *Message, retries *RetryConfig) error {
	var err error
	retry := 0

	for {
		err = currentSender.Send(c, m)

		if err != nil && retries != nil && retry < retries.MaxRetries() {
			backoff := retries.Backoff(retry)

			if retries.ShouldRetry(err) {
				time.Sleep(backoff)
				retry++
				continue
			}
		}
		break
	}

	return err
}

// Sender is anything that can send an email
type Sender interface {
	Send(*Client, *Message) error
}

type defaultSender struct{}

func (s defaultSender) Send(c *Client, m *Message) error {
	return c.Send(m)
}

// DefaultSender is the default SMTP sender
var DefaultSender Sender = defaultSender{}
var currentSender = DefaultSender

// SetSender sets the sender used by Send
func SetSender(sender Sender) {
	currentSender = sender
}
