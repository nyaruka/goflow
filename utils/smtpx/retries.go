package smtpx

import (
	"fmt"
	"strconv"
	"time"
)

// RetryConfig configures if and how retrying of connections happens
type RetryConfig struct {
	Backoffs    []time.Duration
	ShouldRetry func(error) bool
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
func DefaultShouldRetry(err error) bool {
	code := extractCode(err)

	// if error is a transient negative completion reply, it can be retried
	// see https://en.wikipedia.org/wiki/List_of_SMTP_server_return_codes
	return code >= 400 && code < 500
}

// parses an SMTP error response to extract the initial error code
func extractCode(err error) int {
	s := err.Error()

	if len(s) >= 3 {
		code, err := strconv.Atoi(s[0:3])
		if err == nil {
			return code
		}
	}

	return 0
}
