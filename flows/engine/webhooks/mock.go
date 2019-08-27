package webhooks

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/nyaruka/goflow/flows"
)

type mockService struct {
	statusCode int
	body       string
}

// NewMockService creates a new mock webhook service for testing
func NewMockService(statusCode int, body string) flows.WebhookService {
	return &mockService{
		statusCode: statusCode,
		body:       body,
	}
}

func (s *mockService) Call(request *http.Request, resthook string) (*flows.WebhookCall, error) {
	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	recorder.WriteString(s.body)
	recorder.Code = s.statusCode

	response := recorder.Result()
	response.Request = request

	responseTrace, err := httputil.DumpResponse(response, false)
	if err != nil {
		return nil, err
	}

	return &flows.WebhookCall{
		URL:        request.URL.String(),
		Method:     request.Method,
		StatusCode: response.StatusCode,
		Status:     statusFromCode(response.StatusCode, resthook != ""),
		Request:    dump,
		Response:   responseTrace,
		TimeTaken:  1,
		Resthook:   resthook,
	}, nil
}

var _ flows.WebhookService = (*mockService)(nil)
