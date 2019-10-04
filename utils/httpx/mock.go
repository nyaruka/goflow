package httpx

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type mockRequestor struct {
	mocks map[string][]*http.Response
}

func NewMockRequestor(mocks map[string][]*http.Response) Requestor {
	return &mockRequestor{mocks: mocks}
}

func (r *mockRequestor) Do(client *http.Client, request *http.Request) (*http.Response, error) {
	url := request.URL.String()
	mockedResponses := r.mocks[url]
	if len(mockedResponses) == 0 {
		return nil, errors.Errorf("missing mock for URL %s", url)
	}
	mocked := mockedResponses[0]
	r.mocks[url] = mockedResponses[1:]

	mocked.Request = request
	return mocked, nil
}

// NewMockResponse creates a new mock response
func NewMockResponse(status int, body string) *http.Response {
	return &http.Response{
		Status:        fmt.Sprintf("%d OK", status),
		StatusCode:    status,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        nil,
		Body:          ioutil.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
	}
}
