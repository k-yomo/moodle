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
	Token       string
	SecretToken string
}

func Login(ctx context.Context, client *http.Client, serviceURL *url.URL, params *LoginParams) (*LoginResponse, error) {
	serviceURL.Path = path.Join(serviceURL.Path, "/login/token.php")
	urlutil.SetQueries(serviceURL, map[string]string{
		"username": params.Username,
		"password": params.Password,
	})
	req, err := http.NewRequestWithContext(ctx, "GET", serviceURL.String(), nil)
	if err != nil {
		return nil, err
	}

	res := LoginResponse{}
	if err := doAndMap(client, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
