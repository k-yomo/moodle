package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"net/url"
	"path"
)

type AuthAPI interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
}

type authAPI struct {
	httpClient *http.Client
	serviceURL *url.URL
	apiURL     *url.URL
}

func newAuthAPI(httpClient *http.Client, serviceURL, apiURL *url.URL) *authAPI {
	return &authAPI{
		httpClient: httpClient,
		serviceURL: serviceURL,
		apiURL:     apiURL,
	}
}

type LoginResponse struct {
	Token        string `json:"token"`
	PrivateToken string `json:"privatetoken"`
}

func (a *authAPI) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	u := urlutil.CopyWithQueries(a.serviceURL, map[string]string{
		"username": username,
		"password": password,
		"service":  "moodle_mobile_app",
	})
	u.Path = path.Join(u.Path, "/login/token.php")
	res := LoginResponse{}
	if err := getAndUnmarshal(ctx, a.httpClient, u, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
