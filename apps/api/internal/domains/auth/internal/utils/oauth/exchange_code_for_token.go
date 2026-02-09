package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func ExchangeCodeForToken(
	ctx context.Context,
	httpClient *http.Client,
	tokenURL string,
	clientID string,
	clientSecret string,
	redirectURI string,
	code string,
) (TokenResponse, error) {
	reqBody := tokenRequest{
		GrantType:    "authorization_code",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Code:         code,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewReader(body))
	if err != nil {
		return TokenResponse{}, fmt.Errorf("create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	bodyBytes, statusCode, err := httpapi.DoRequest(
		ctx,
		httpClient,
		req.Method,
		req.URL.String(),
		req.Body,
		map[string]string{
			"Content-Type": req.Header.Get("Content-Type"),
			"Accept":       req.Header.Get("Accept"),
		},
	)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("do token request: %w", err)
	}

	if statusCode < 200 || statusCode >= 300 {
		return TokenResponse{}, fmt.Errorf("token request failed: status=%d body=%s", statusCode, string(bodyBytes))
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		return TokenResponse{}, fmt.Errorf("unmarshal token response: %w", err)
	}

	if tokenResponse.AccessToken == "" {
		return TokenResponse{}, fmt.Errorf("token response missing access_token")
	}

	return tokenResponse, nil
}
