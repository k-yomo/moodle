package moodle

import (
	"net/http"
	"net/url"
)

type UserAPI interface {
}

type userAPI struct {
	httpClient *http.Client
	apiURL     *url.URL
}

func newUserAPI(httpClient *http.Client, apiURL *url.URL) *userAPI {
	return &userAPI{
		httpClient: httpClient,
		apiURL:     apiURL,
	}
}
