package httpx

import (
	"net/http"
	"strconv"
	"time"

	"github.com/nyaruka/goflow/utils/dates"
)

// RetryConfig configures if and how retrying of requests happens
type RetryConfig struct {
	Delays      []time.Duration
	ShouldRetry func(*http.Request, *http.Response, time.Duration) bool
}

// NewRetryDelays creates a new retry config with the given delays
func NewRetryDelays(delays ...time.Duration) *RetryConfig {
	return &RetryConfig{Delays: delays, ShouldRetry: DefaultShouldRetry}
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
