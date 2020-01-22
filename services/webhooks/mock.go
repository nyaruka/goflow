package webhooks

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/nyaruka/goflow/flows"
)

type mockService struct {
	statusCode  int
	contentType string
	body        string
}

// NewMockService creates a new mock webhook service for testing
func NewMockService(statusCode int, contentType, body string) flows.WebhookService {
	return &mockService{
		statusCode:  statusCode,
		contentType: contentType,
		body:        body,
	}
}

func (s *mockService) Call(session flows.Session, request *http.Request) (*flows.WebhookCall, error) {
	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", s.contentType)
	recorder.WriteString(s.body)
	recorder.Code = s.statusCode

	response := recorder.Result()
	response.Request = request

	responseTrace, err := httputil.DumpResponse(response, true)
	if err != nil {
		return nil, err
	}

	return &flows.WebhookCall{
		URL:        request.URL.String(),
		Method:     request.Method,
		StatusCode: response.StatusCode,
		Request:    dump,
		Response:   responseTrace,
		TimeTaken:  1,
	}, nil
}

var _ flows.WebhookService = (*mockService)(nil)
