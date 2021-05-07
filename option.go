package moodle

import "net/http"

type ClientOptions struct {
	Token      string
	HttpClient *http.Client
	Debug      bool
}

func newDefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		HttpClient: http.DefaultClient,
		Debug:      false,
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

// WithDebugEnabled enable debug logs
// this option is should be used in development only.
func WithDebugEnabled() ClientOption {
	return newClientOptionFunc(func(c *ClientOptions) {
		c.Debug = true
	})
}

func withToken(token string) ClientOption {
	return newClientOptionFunc(func(c *ClientOptions) {
		c.Token = token
	})
}
