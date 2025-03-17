package testutil

import (
	"bytes"
	"io"
	"net/http"
)

// MockHTTPClient is a mock implementation of an HTTP client for testing.
type MockHTTPClient struct {
	Requests *http.Request
	Response []byte
}

// Do implements the HTTPClient interface by recording the request and returning the mock response.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.Requests = req
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer(m.Response)),
	}, nil
}
