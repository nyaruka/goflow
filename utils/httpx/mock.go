package httpx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/nyaruka/goflow/utils/dates"
	"github.com/pkg/errors"
)

type MockRequestor struct {
	mocks map[string][]MockResponse
}

func NewMockRequestor(mocks map[string][]MockResponse) *MockRequestor {
	return &MockRequestor{mocks: mocks}
}

func (r *MockRequestor) Do(client *http.Client, request *http.Request) (*http.Response, error) {
	url := request.URL.String()
	mockedResponses := r.mocks[url]
	if len(mockedResponses) == 0 {
		return nil, errors.Errorf("missing mock for URL %s", url)
	}

	// pop the next mocked response for this URL
	mocked := mockedResponses[0]
	r.mocks[url] = mockedResponses[1:]

	return mocked.Make(request), nil
}

func (r *MockRequestor) HasUnused() bool {
	for _, mocks := range r.mocks {
		if len(mocks) > 0 {
			return true
		}
	}
	return false
}

func (r *MockRequestor) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.mocks)
}

var _ Requestor = (*MockRequestor)(nil)

type MockResponse struct {
	Status int    `json:"status" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

func (m MockResponse) Make(request *http.Request) *http.Response {
	return &http.Response{
		Request:       request,
		Status:        fmt.Sprintf("%d OK", m.Status),
		StatusCode:    m.Status,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        nil,
		Body:          ioutil.NopCloser(strings.NewReader(m.Body)),
		ContentLength: int64(len(m.Body)),
	}
}

// NewMockResponse creates a new mock response
func NewMockResponse(status int, body string) MockResponse {
	return MockResponse{status, body}
}

// NewMockTrace creates a new trace for testing without making an actual request
func NewMockTrace(method, url string, status int, body string) *Trace {
	request, _ := http.NewRequest(method, url, nil)
	requestTrace, _ := httputil.DumpRequestOut(request, true)

	response := NewMockResponse(status, body).Make(request)
	responseTrace, _ := httputil.DumpResponse(response, true)

	return &Trace{
		Request:       request,
		Response:      response,
		RequestTrace:  requestTrace,
		ResponseTrace: responseTrace,
		Body:          []byte(body),
		StartTime:     dates.Now(),
		EndTime:       dates.Now(),
	}
}
