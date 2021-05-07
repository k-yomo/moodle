package moodle

import (
	"net/http"
	"testing"
)

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	clientOptions := ClientOptions{}
	httpClient := &http.Client{}

	WithHTTPClient(httpClient).apply(&clientOptions)
	if clientOptions.HttpClient != httpClient {
		t.Errorf("WithHTTPClient() = %v, want %v", clientOptions.HttpClient, httpClient)
	}
}

func TestWithDebugEnabled(t *testing.T) {
	t.Parallel()

	clientOptions := ClientOptions{}
	WithDebugEnabled().apply(&clientOptions)
	if clientOptions.Debug != true {
		t.Errorf("WithDebugEnabled() = %v, want %v", clientOptions.Debug, true)
	}
}
