package smtpx

import (
	"fmt"
	"strings"
	"time"
)

// RetryConfig configures if and how retrying of connections happens
type RetryConfig struct {
	Backoffs    []time.Duration
	ShouldRetry func(bool, error) bool
}

// NewFixedRetries creates a new retry config with the given backoffs
func NewFixedRetries(backoffs ...time.Duration) *RetryConfig {
	return &RetryConfig{Backoffs: backoffs, ShouldRetry: DefaultShouldRetry}
}

// MaxRetries gets the maximum number of retries allowed
func (r *RetryConfig) MaxRetries() int {
	return len(r.Backoffs)
}

// Backoff gets the backoff time for the nth retry
func (r *RetryConfig) Backoff(n int) time.Duration {
	if n >= len(r.Backoffs) {
		panic(fmt.Sprintf("%d not a valid retry number for this config", n))
	}
	return r.Backoffs[n]
}

// DefaultShouldRetry is the default function for determining if a send should be retried
func DefaultShouldRetry(fromDial bool, err error) bool {
	errMsg := strings.ToLower(err.Error())

	// if the error message looks to be an authentication failure, don't retry
	if strings.Contains(errMsg, "username") || strings.Contains(errMsg, "password") {
		return false
	}

	// otherwise retry if the err came from the dial stage
	return fromDial
}

// Send an email using SMTP
func Send(c *Client, m *Message, retries *RetryConfig) error {
	var fromDial bool
	var err error
	retry := 0

	for {
		fromDial, err = currentSender.Send(c, m)

		if retries != nil && retry < retries.MaxRetries() {
			backoff := retries.Backoff(retry)

			if retries.ShouldRetry(fromDial, err) {
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
	Send(*Client, *Message) (bool, error)
}

type defaultSender struct{}

func (s defaultSender) Send(c *Client, m *Message) (bool, error) {
	return c.Send(m)
}

// DefaultSender is the default SMTP sender
var DefaultSender Sender = defaultSender{}
var currentSender = DefaultSender

// SetSender sets the sender used by Send
func SetSender(sender Sender) {
	currentSender = sender
}
