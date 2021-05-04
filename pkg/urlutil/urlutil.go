package urlutil

import "net/url"

func Copy(u *url.URL) *url.URL {
	copied := *u
	return &copied
}

func SetQueries(u *url.URL, queries ...map[string]string) {
	q := u.Query()
	for _, query := range queries {
		for k, v := range query {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()
}

func CopyWithQueries(original *url.URL, queries ...map[string]string) *url.URL {
	u := Copy(original)
	SetQueries(u, queries...)
	return u
}
