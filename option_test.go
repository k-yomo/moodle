package moodle

import (
	"net/http"
	"testing"
)

func TestWithHTTPClient(t *testing.T) {

	clientOptions := ClientOptions{}
	httpClient := &http.Client{}

	WithHTTPClient(httpClient).apply(&clientOptions)
	if clientOptions.HttpClient != httpClient {
		t.Errorf("WithHTTPClient() = %v, want %v", clientOptions.HttpClient, httpClient)
	}
}
