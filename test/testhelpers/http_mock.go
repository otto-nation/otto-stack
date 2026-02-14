package testhelpers

import (
	"io"
	"net/http"
	"strings"
)

// MockHTTPClient is a mock implementation of http.Client for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("{}")),
	}, nil
}

// MockResponse creates a mock HTTP response
func MockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

// MockJSONResponse creates a mock HTTP response with JSON content type
func MockJSONResponse(statusCode int, body string) *http.Response {
	resp := MockResponse(statusCode, body)
	resp.Header.Set("Content-Type", "application/json")
	return resp
}
