package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"net/url"
	"path"
)

type LoginParams struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token        string `json:"token"`
	PrivateToken string `json:"privatetoken"`
}

func Login(ctx context.Context, client *http.Client, serviceURL *url.URL, params *LoginParams) (*LoginResponse, error) {
	u := urlutil.Copy(serviceURL)
	u.Path = path.Join(u.Path, "/login/token.php")
	urlutil.SetQueries(u, map[string]string{
		"username": params.Username,
		"password": params.Password,
		"service":  "moodle_mobile_app",
	})
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	res := LoginResponse{}
	if err := doAndMap(client, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
