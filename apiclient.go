package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"net/url"
)

type apiClient struct {
	httpClient *http.Client
	apiURL     *url.URL
}

func newAPIClient(httpClient *http.Client, apiURL *url.URL) *apiClient {
	return &apiClient{
		httpClient: httpClient,
		apiURL:     apiURL,
	}
}

// callMoodleFunction call moodle's service function and map the response json to `to` param.
func (a *apiClient) callMoodleFunction(ctx context.Context, to interface{}, queryParams ...map[string]string) error {
	u := urlutil.CopyWithQueries(a.apiURL, queryParams...)
	return getAndUnmarshal(ctx, a.httpClient, u, to)
}
