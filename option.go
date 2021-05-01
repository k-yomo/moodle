package moodle

import "net/http"

type ClientOptions struct {
	Token      string
	HttpClient *http.Client
}

func newDefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		HttpClient: http.DefaultClient,
	}
}

// ClientOption is a option to change client configuration.
type ClientOption interface {
	apply(*ClientOptions)
}

type clientOptionFunc struct {
	f func(config *ClientOptions)
}

func (c *clientOptionFunc) apply(pc *ClientOptions) {
	c.f(pc)
}

func newClientOptionFunc(f func(pc *ClientOptions)) *clientOptionFunc {
	return &clientOptionFunc{
		f: f,
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return newClientOptionFunc(func(c *ClientOptions) {
		c.HttpClient = httpClient
	})
}

func withToken(token string) ClientOption {
	return newClientOptionFunc(func(c *ClientOptions) {
		c.Token = token
	})
}
