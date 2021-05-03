package moodle

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	serviceURL, _ := url.Parse("https://test.edu")
	got, err := NewClient(context.Background(), serviceURL, "test")
	if err != nil {
		t.Fatalf("NewClientWithLogin() error = %v", err)
	}

	if u := got.serviceURL.String(); u != "https://test.edu" {
		t.Errorf("NewClientWithLogin(), got.serviceURL = %v, want = %v", u, "https://test.edu")
	}
	if u := got.apiURL.String(); u != "https://test.edu/webservice/rest/server.php?moodlewsrestformat=json&wstoken=test" {
		t.Errorf("NewClientWithLogin(), got.apiURL = %v, want = %v", u, "https://test.edu/webservice/rest/server.php?moodlewsrestformat=json&wstoken=test")
	}
	if got.UserAPI == nil {
		t.Errorf("NewClientWithLogin(), got.UserAPI = nil")
	}
	if got.CourseAPI == nil {
		t.Errorf("NewClientWithLogin(), got.CourseAPI = nil")
	}
	if got.QuizAPI == nil {
		t.Errorf("NewClientWithLogin(), got.QuizAPI = nil")
	}
}

func TestNewClientWithLogin(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"token":"test", "privatetoken": "private"}`)
	})
	s := httptest.NewServer(h)
	serviceURL, _ := url.Parse(s.URL)

	got, err := NewClientWithLogin(context.Background(), serviceURL, "", "")
	if err != nil {
		t.Fatalf("NewClientWithLogin() error = %v", err)
	}

	if u := got.serviceURL.String(); u != serviceURL.String() {
		t.Errorf("NewClientWithLogin(), got.serviceURL = %v, want = %v", u, serviceURL.String())
	}
	if u := got.apiURL.String(); u != serviceURL.String()+"/webservice/rest/server.php?moodlewsrestformat=json&wstoken=test" {
		t.Errorf("NewClientWithLogin(), got.apiURL = %v, want = %v", u, serviceURL.String()+"/webservice/rest/server.php?moodlewsrestformat=json&wstoken=test")
	}
	if got.CourseAPI == nil {
		t.Errorf("NewClientWithLogin(), got.CourseAPI = nil")
	}
	if got.QuizAPI == nil {
		t.Errorf("NewClientWithLogin(), got.QuizAPI = nil")
	}
}
