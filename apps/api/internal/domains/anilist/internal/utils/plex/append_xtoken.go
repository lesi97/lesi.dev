package plex

import "net/url"

func (s *Store) appendXToken(path string) (*string, error) {
	url, err := url.ParseRequestURI(s.env.BaseUrl + path)
	if err != nil {
		return nil, err
	}
	query := url.Query()
	query.Add("X-Plex-Token", s.env.XToken)
	url.RawQuery = query.Encode()
	str := url.String()
	return &str, nil
}
