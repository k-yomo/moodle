package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"net/url"
	"path"
)

// Client is a Moodle API client scoped to a service.
type Client struct {
	serviceURL *url.URL
	apiURL     *url.URL
	opts       *ClientOptions

	AuthAPI   AuthAPI
	CourseAPI CourseAPI
	QuizAPI   QuizAPI
}

// NewClient creates a new Moodle client.
func NewClient(ctx context.Context, serviceURL *url.URL, token string, opt ...ClientOption) (*Client, error) {
	// TODO: Validate token
	return newClient(serviceURL, append(opt, withToken(token))...), nil
}

// NewClientWithLogin creates a new Moodle client with token retrieved from login request.
func NewClientWithLogin(ctx context.Context, serviceURL *url.URL, username, password string, opt ...ClientOption) (*Client, error) {
	res, err := newAuthAPI(http.DefaultClient, serviceURL).Login(ctx, username, password)
	if err != nil {
		return nil, err
	}

	c := newClient(serviceURL, append(opt, withToken(res.Token))...)
	return c, nil
}

func newClient(serviceURL *url.URL, opt ...ClientOption) *Client {
	opts := newDefaultClientOptions()
	for _, o := range opt {
		o.apply(opts)
	}

	apiURL := *serviceURL
	apiURL.Path = path.Join(apiURL.Path, "/webservice/rest/server.php")
	urlutil.SetQueries(&apiURL, map[string]string{
		"moodlewsrestformat": "json",
		"wstoken":            opts.Token,
	})

	return &Client{
		serviceURL: serviceURL,
		apiURL:     &apiURL,
		opts:       opts,

		CourseAPI: newCourseAPI(opts.HttpClient, &apiURL),
		QuizAPI:   newQuizAPI(opts.HttpClient, &apiURL),
	}
}
