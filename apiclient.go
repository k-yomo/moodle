package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

type apiClient struct {
	authToken  string
	httpClient *http.Client
	serviceURL *url.URL
	apiURL     *url.URL
	debug      bool
}

func newAPIClient(httpClient *http.Client, serviceURL *url.URL, authToken string, debug bool) *apiClient {
	apiURL := urlutil.Copy(serviceURL)
	apiURL.Path = path.Join(apiURL.Path, "/webservice/rest/server.php")
	urlutil.SetQueries(apiURL, map[string]string{
		"moodlewsrestformat": "json",
		"wstoken":            authToken,
	})

	return &apiClient{
		authToken:  authToken,
		httpClient: httpClient,
		serviceURL: serviceURL,
		apiURL:     apiURL,
		debug:      debug,
	}
}

func (a *apiClient) updateToken(authToken string) {
	a.authToken = authToken
	urlutil.SetQueries(a.apiURL, map[string]string{
		"wstoken": authToken,
	})
}

// callMoodleFunction call moodle's service function and map the response json to `to` param.
func (a *apiClient) callMoodleFunction(ctx context.Context, to interface{}, queryParams ...map[string]string) error {
	u := urlutil.CopyWithQueries(a.apiURL, queryParams...)
	return a.getAndUnmarshal(ctx, u, to)
}

func (a *apiClient) getAndUnmarshal(ctx context.Context, u *url.URL, to interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}
	if a.debug {
		log.Printf(`[INFO] make http request
	method: %s
	url: %s
`, req.Method, req.URL.String())
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if a.debug {
		log.Printf(`[INFO] raw response
	status: %s
	body: %s

`, resp.Status, string(bodyBytes))
	}
	if err := mapResponseBodyToStruct(bodyBytes, to); err != nil {
		return err
	}
	return nil
}
