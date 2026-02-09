package twitch_frontend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/httpapi"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func ExchangeCodeForToken(
	ctx context.Context,
	tokenBaseURL string,
	clientID string,
	clientSecret string,
	redirectURL string,
	code string,
	pkceVerifier string,
) (string, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}

	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", redirectURL)
	form.Set("code_verifier", pkceVerifier)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%v/token", tokenBaseURL),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	bodyBytes, statusCode, err := httpapi.DoRequest(
		ctx,
		httpClient,
		req.Method,
		req.URL.String(),
		req.Body,
		map[string]string{
			"Content-Type": req.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", err
	}

	if statusCode < 200 || statusCode >= 300 {
		return "", fmt.Errorf("twitch token exchange failed")
	}

	var tr tokenResponse
	if err := json.Unmarshal(bodyBytes, &tr); err != nil {
		return "", err
	}

	if tr.AccessToken == "" {
		return "", fmt.Errorf("missing access token")
	}

	return tr.AccessToken, nil
}
