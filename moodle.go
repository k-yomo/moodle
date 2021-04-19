package moodle

import (
	"context"
	"github.com/k-yomo/moodle/moodleapi"
	"net/url"
)

// Client is a Moodle API client scoped to a service.
type Client struct {
	serviceURL *url.URL
	opts       *ClientOptions
}

// NewClient creates a new Moodle client.
func NewClient(ctx context.Context, serviceURL *url.URL, token string, opt ...ClientOption) (*Client, error) {
	// TODO: Validate token
	return newClient(serviceURL, append(opt, withToken(token))...), nil
}

// NewClientWithLogin creates a new Moodle client with token retrieved from login request.
func NewClientWithLogin(ctx context.Context, serviceURL *url.URL, loginParams *moodleapi.LoginParams, opt ...ClientOption) (*Client, error) {
	c := newClient(serviceURL, opt...)
	resp, err := moodleapi.Login(
		ctx,
		c.opts.HttpClient,
		serviceURL,
		loginParams,
	)
	if err != nil {
		return nil, err
	}
	withToken(resp.Token).apply(c.opts)
	return c, nil
}

func newClient(serviceURL *url.URL, opt ...ClientOption) *Client {
	opts := newDefaultClientOptions()
	for _, o := range opt {
		o.apply(opts)
	}
	return &Client{
		serviceURL: serviceURL,
		opts:       opts,
	}
}
