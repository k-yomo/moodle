package moodle

import (
	"context"
	"net/url"
)

// Client is a Moodle API client scoped to a service.
type Client struct {
	opts      *ClientOptions
	apiClient *apiClient

	AuthAPI   AuthAPI
	SiteAPI   SiteAPI
	UserAPI   UserAPI
	CourseAPI CourseAPI
	QuizAPI   QuizAPI
	GradeAPI  GradeAPI
}

// NewClient creates a new Moodle client.
func NewClient(ctx context.Context, serviceURL *url.URL, token string, opt ...ClientOption) (*Client, error) {
	// TODO: Validate token
	return newClient(serviceURL, append(opt, withToken(token))...), nil
}

// NewClientWithLogin creates a new Moodle client with token retrieved from login request.
func NewClientWithLogin(ctx context.Context, serviceURL *url.URL, username, password string, opt ...ClientOption) (*Client, error) {
	c := newClient(serviceURL, opt...)
	res, err := c.AuthAPI.Login(ctx, username, password)
	if err != nil {
		return nil, err
	}
	c.apiClient.updateToken(res.Token)

	return c, nil
}

func newClient(serviceURL *url.URL, opt ...ClientOption) *Client {
	opts := newDefaultClientOptions()
	for _, o := range opt {
		o.apply(opts)
	}
	apiClient := newAPIClient(opts.HttpClient, serviceURL, opts.Token, opts.Debug)

	return &Client{
		opts:      opts,
		apiClient: apiClient,
		AuthAPI:   newAuthAPI(apiClient),
		SiteAPI:   newSiteAPI(apiClient),
		UserAPI:   newUserAPI(apiClient),
		CourseAPI: newCourseAPI(apiClient),
		QuizAPI:   newQuizAPI(apiClient),
		GradeAPI:  newGradeAPI(apiClient),
	}
}
