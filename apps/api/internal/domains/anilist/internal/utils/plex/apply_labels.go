package plex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func (s *Store) ApplyPlexLabels(ctx context.Context, ratingKey string, tagType string, labels []string) error {
	base, err := s.appendXToken(fmt.Sprintf("/library/metadata/%s", ratingKey))
	if err != nil {
		return err
	}
	u, err := url.Parse(*base)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set(fmt.Sprintf("%s.locked", tagType), "1")

	for i, lab := range labels {
		q.Set(fmt.Sprintf("%s[%d].tag.tag", tagType, i), lab)
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	_, statusCode, err := httpapi.DoRequest(ctx, client, req.Method, req.URL.String(), req.Body, nil)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return fmt.Errorf("plex returned status %d", statusCode)
	}

	return nil
}
