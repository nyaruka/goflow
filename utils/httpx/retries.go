package httpx

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/random"
)

// RetryConfig configures if and how retrying of requests happens
type RetryConfig struct {
	Backoffs    []time.Duration
	Jitter      float64
	ShouldRetry func(*http.Request, *http.Response, time.Duration) bool
}

// NewFixedRetries creates a new retry config with the given backoffs
func NewFixedRetries(backoffs ...time.Duration) *RetryConfig {
	return &RetryConfig{Backoffs: backoffs, ShouldRetry: DefaultShouldRetry}
}

// NewExponentialRetries creates a new retry config with the given delays
func NewExponentialRetries(initialBackoff time.Duration, count int, jitter float64) *RetryConfig {
	backoffs := make([]time.Duration, count)
	backoffs[0] = initialBackoff
	for i := 1; i < count; i++ {
		backoffs[i] = backoffs[i-1] * 2
	}

	return &RetryConfig{Backoffs: backoffs, Jitter: jitter, ShouldRetry: DefaultShouldRetry}
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

	base := r.Backoffs[n]
	jitter := time.Duration(r.Jitter * float64(random.IntN(int(base))-(int(base)/2)))
	return base + jitter
}

// DefaultShouldRetry is the default function for determining if a response should be retried
func DefaultShouldRetry(request *http.Request, response *http.Response, withDelay time.Duration) bool {
	// any response with a Retry-After header is candidate for a retry (usually used with 301, 429, 503 status codes)
	if response != nil {
		retryAfter := response.Header.Get("Retry-After")
		if retryAfter != "" {
			requestedDelay := ParseRetryAfter(retryAfter)

			// as long as the server has requested a delay which is less than or equal to what we intended
			if requestedDelay != 0 && requestedDelay <= withDelay {
				return true
			}
		}
	}

	// otherwise retry if request is idempotent and response is a failure (excluding 500 and 501)
	return isIdempotent(request) && (response == nil || response.StatusCode > 501)
}

// see https://github.com/golang/go/blob/100bf440b9a69c6dce8daeebed038d607c963b8f/src/net/http/request.go#L1395
func isIdempotent(r *http.Request) bool {
	switch r.Method {
	case "GET", "HEAD", "OPTIONS", "TRACE":
		return true
	}

	return r.Header.Get("Idempotency-Key") != "" || r.Header.Get("X-Idempotency-Key") != ""
}

// ParseRetryAfter parses value of Retry-After headers which can be date or delay in seconds
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Retry-After
func ParseRetryAfter(value string) time.Duration {
	asTime, err := http.ParseTime(value)
	if err == nil {
		delta := asTime.Sub(dates.Now())
		if delta >= 0 {
			return delta
		}
	} else {
		asSeconds, err := strconv.Atoi(value)
		if err == nil {
			return time.Duration(asSeconds) * time.Second
		}
	}

	return 0
}
