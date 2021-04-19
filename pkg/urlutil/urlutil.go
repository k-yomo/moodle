package urlutil

import "net/url"

func SetQueries(u *url.URL, queries map[string]string)  {
	q := u.Query()
	for k, v := range queries {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
}
