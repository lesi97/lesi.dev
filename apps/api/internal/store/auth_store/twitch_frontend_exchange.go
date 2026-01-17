package auth_store

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
)

type TwitchFrontendToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (s *AuthStore) twitchFrontendExchange(
	ctx context.Context,
	code string,
	pkceVerifier string,
) (*TwitchFrontendToken, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}

	form := url.Values{}
	form.Set("client_id", *s.twitch.client_id)
	form.Set("client_secret", *s.twitch.client_secret)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", fmt.Sprintf("%v/api/auth/v1/callback", *s.webUrl))
	form.Set("code_verifier", pkceVerifier)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%v/token", s.twitch.base_url),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}
		fmt.Printf("Response body:\n")
		utils.PrintPrettyJSON(bodyBytes)
		return nil, fmt.Errorf("twitch token exchange failed")
	}

	var tr TwitchFrontendToken
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}

	if tr.AccessToken == "" {
		return nil, fmt.Errorf("missing access token")
	}

	return &tr, nil
}
