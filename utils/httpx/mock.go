package httpx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
		panic(fmt.Sprintf("missing mock for URL %s", url))
	}

	// pop the next mocked response for this URL
	mocked := mockedResponses[0]
	r.mocks[url] = mockedResponses[1:]

	if mocked.Status == 0 {
		return nil, errors.New("unable to connect to server")
	}

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
		Status:        fmt.Sprintf("%d %s", m.Status, http.StatusText(m.Status)),
		StatusCode:    m.Status,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        nil,
		Body:          ioutil.NopCloser(strings.NewReader(m.Body)),
		ContentLength: int64(len(m.Body)),
	}
}

// MockConnectionError mocks a connection error
var MockConnectionError = MockResponse{0, ""}

// NewMockResponse creates a new mock response
func NewMockResponse(status int, body string) MockResponse {
	return MockResponse{status, body}
}
