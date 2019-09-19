package webhooks

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/nyaruka/goflow/flows"
)

type mockProvider struct {
	statusCode  int
	contentType string
	body        string
}

// NewMockProvider creates a new mock webhook provider for testing
func NewMockProvider(statusCode int, contentType, body string) flows.WebhookProvider {
	return &mockProvider{
		statusCode:  statusCode,
		contentType: contentType,
		body:        body,
	}
}

func (s *mockProvider) Call(request *http.Request, resthook string) (*flows.WebhookCall, error) {
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
		Status:     statusFromCode(response.StatusCode, resthook != ""),
		Request:    dump,
		Response:   responseTrace,
		TimeTaken:  1,
		Resthook:   resthook,
	}, nil
}

var _ flows.WebhookProvider = (*mockProvider)(nil)
