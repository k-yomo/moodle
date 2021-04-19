package urlutil

import "net/url"

func Copy(u *url.URL) *url.URL {
	copied := *u
	return &copied
}

func SetQueries(u *url.URL, queries map[string]string) {
	q := u.Query()
	for k, v := range queries {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
}
