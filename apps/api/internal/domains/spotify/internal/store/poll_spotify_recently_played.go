package store

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

const (
	spotifyApplication = "Spotify"
	spotifyPollSource  = "spotify-poll"
	spotifyAccountsURL = "https://accounts.spotify.com"
	spotifyAPIURL      = "https://api.spotify.com/v1"
)

type spotifyAPIAuth struct {
	AccessToken string
	APIURL      string
	Details     *db.ApiDetails
}

type spotifyTokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type spotifyRecentlyPlayedResponse struct {
	Items []json.RawMessage `json:"items"`
}

type spotifyPlayHistory struct {
	Track    spotifyTrack `json:"track"`
	PlayedAt time.Time    `json:"played_at"`
}

type spotifyTrack struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	ExternalURLs spotifyExternalURLs `json:"external_urls"`
	Album        spotifyAlbum        `json:"album"`
	Artists      []spotifyArtist     `json:"artists"`
}

type spotifyAlbum struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	ExternalURLs spotifyExternalURLs `json:"external_urls"`
	Images       []spotifyImage      `json:"images"`
}

type spotifyArtistsResponse struct {
	Artists []spotifyArtist `json:"artists"`
}

type spotifyArtist struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	ExternalURLs spotifyExternalURLs `json:"external_urls"`
	Images       []spotifyImage      `json:"images"`
}

type spotifyExternalURLs struct {
	Spotify string `json:"spotify"`
}

type spotifyImage struct {
	URL string `json:"url"`
}

func (s *Store) PollSpotifyRecentlyPlayed(ctx context.Context, input model.SpotifyPollInput) (*model.SpotifyPollResult, error) {
	limit := input.Limit
	if limit == 0 {
		limit = 50
	}
	if limit < 1 || limit > 50 {
		return nil, fmt.Errorf("limit must be between 1 and 50")
	}

	auth, err := s.spotifyAPIAuth(ctx)
	if err != nil {
		return nil, err
	}
	s.refreshCurrentlyPlayingCacheAsync(ctx, auth)

	afterMS := input.AfterMS
	if afterMS == nil {
		afterMS, err = s.latestSpotifyScrobbleMS(ctx)
		if err != nil {
			return nil, err
		}
	}

	recentlyPlayed, err := s.fetchSpotifyRecentlyPlayed(ctx, auth.APIURL, auth.AccessToken, limit, afterMS)
	if errors.Is(err, errSpotifyUnauthorized) {
		accessToken, refreshErr := s.refreshSpotifyAccessToken(ctx, auth.Details)
		if refreshErr != nil {
			return nil, refreshErr
		}
		auth.AccessToken = accessToken
		recentlyPlayed, err = s.fetchSpotifyRecentlyPlayed(ctx, auth.APIURL, auth.AccessToken, limit, afterMS)
	}
	if err != nil {
		return nil, err
	}

	artistImages := s.fetchSpotifyArtistImages(ctx, auth.APIURL, auth.AccessToken, recentlyPlayed)
	result := &model.SpotifyPollResult{
		Fetched: len(recentlyPlayed.Items),
		AfterMS: afterMS,
		Source:  spotifyPollSource,
	}
	var latestPlayedAt time.Time

	for _, rawItem := range recentlyPlayed.Items {
		var item spotifyPlayHistory
		if err := json.Unmarshal(rawItem, &item); err != nil {
			result.Skipped++
			continue
		}

		input, ok, err := spotifyPlayHistoryToScrobbleInput(item, rawItem, artistImages)
		if err != nil {
			return nil, err
		}
		if !ok {
			result.Skipped++
			continue
		}

		if _, err := s.InsertScrobble(ctx, input); err != nil {
			return nil, err
		}

		result.Scrobbled++
		if latestPlayedAt.IsZero() || item.PlayedAt.After(latestPlayedAt) {
			latestPlayedAt = item.PlayedAt
			playedAt := item.PlayedAt.UTC().Format(time.RFC3339)
			result.LatestPlayedAt = &playedAt
		}
	}

	return result, nil
}

var errSpotifyUnauthorized = errors.New("spotify unauthorized")

func (s *Store) spotifyAPIAuth(ctx context.Context) (spotifyAPIAuth, error) {
	if s.DB == nil {
		return spotifyAPIAuth{}, fmt.Errorf("music store database is not configured")
	}

	apiDetails, err := s.DB.FetchApiDetails(ctx, spotifyApplication, s.Logger)
	if err != nil {
		return spotifyAPIAuth{}, err
	}
	if apiDetails == nil {
		return spotifyAPIAuth{}, fmt.Errorf("missing Spotify API details; authenticate Spotify first")
	}

	apiURL := spotifyAPIURL
	if apiDetails.BaseURL != nil && strings.TrimSpace(*apiDetails.BaseURL) != "" {
		apiURL = strings.TrimRight(strings.TrimSpace(*apiDetails.BaseURL), "/")
	}

	accessToken := stringValue(apiDetails.AccessToken)
	if accessToken != "" && !spotifyAccessTokenExpired(apiDetails.RefreshTokenExpiry) {
		return spotifyAPIAuth{
			AccessToken: accessToken,
			APIURL:      apiURL,
			Details:     apiDetails,
		}, nil
	}

	accessToken, err = s.refreshSpotifyAccessToken(ctx, apiDetails)
	if err != nil {
		return spotifyAPIAuth{}, err
	}

	return spotifyAPIAuth{
		AccessToken: accessToken,
		APIURL:      apiURL,
		Details:     apiDetails,
	}, nil
}

func (s *Store) refreshSpotifyAccessToken(ctx context.Context, apiDetails *db.ApiDetails) (string, error) {
	clientID := stringValue(apiDetails.ClientID)
	if clientID == "" {
		clientID = strings.TrimSpace(os.Getenv("SPOTIFY_CLIENT_ID"))
	}
	clientSecret := stringValue(apiDetails.ClientSecret)
	if clientSecret == "" {
		clientSecret = strings.TrimSpace(os.Getenv("SPOTIFY_CLIENT_SECRET"))
	}
	refreshToken := stringValue(apiDetails.RefreshToken)

	if clientID == "" {
		return "", fmt.Errorf("missing Spotify client id")
	}
	if clientSecret == "" {
		return "", fmt.Errorf("missing Spotify client secret")
	}
	if refreshToken == "" {
		return "", fmt.Errorf("missing Spotify refresh token; authenticate Spotify first")
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	tokenURL := spotifyAccountsURL + "/api/token"
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret))
	body, statusCode, err := httpapi.DoRequest(
		ctx,
		s.httpClient(),
		http.MethodPost,
		tokenURL,
		strings.NewReader(form.Encode()),
		map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/x-www-form-urlencoded",
			"Accept":        "application/json",
		},
	)
	if err != nil {
		return "", err
	}
	if statusCode < 200 || statusCode >= 300 {
		return "", fmt.Errorf("spotify token refresh failed: status=%d body=%s", statusCode, string(body))
	}

	var tokenResp spotifyTokenRefreshResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}
	if strings.TrimSpace(tokenResp.AccessToken) == "" {
		return "", fmt.Errorf("spotify token refresh response missing access_token")
	}
	if tokenResp.ExpiresIn <= 0 {
		tokenResp.ExpiresIn = 3600
	}

	accessToken := tokenResp.AccessToken
	var refreshTokenPtr *string
	if strings.TrimSpace(tokenResp.RefreshToken) != "" {
		refreshToken = tokenResp.RefreshToken
		refreshTokenPtr = &refreshToken
		apiDetails.RefreshToken = refreshTokenPtr
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UnixMilli()
	apiDetails.AccessToken = &accessToken
	apiDetails.RefreshTokenExpiry = &expiresAt

	if err := s.DB.UpdateApiDetails(
		ctx,
		spotifyApplication,
		nil,
		nil,
		&accessToken,
		refreshTokenPtr,
		&expiresAt,
		apiDetails.BaseURL,
		apiDetails.RedirectURL,
		s.Logger,
	); err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *Store) latestSpotifyScrobbleMS(ctx context.Context) (*int64, error) {
	var scrobbledAt time.Time
	err := s.DB.QueryRow(ctx, `
		SELECT scrobbled_at
		FROM music.scrobbles
		WHERE source = $1
		ORDER BY scrobbled_at DESC
		LIMIT 1
	`, spotifyPollSource).Scan(&scrobbledAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	afterMS := scrobbledAt.UTC().UnixMilli()
	return &afterMS, nil
}

func (s *Store) fetchSpotifyRecentlyPlayed(ctx context.Context, apiURL string, accessToken string, limit int, afterMS *int64) (spotifyRecentlyPlayedResponse, error) {
	endpoint, err := url.Parse(strings.TrimRight(apiURL, "/") + "/me/player/recently-played")
	if err != nil {
		return spotifyRecentlyPlayedResponse{}, err
	}

	query := endpoint.Query()
	query.Set("limit", strconv.Itoa(limit))
	if afterMS != nil {
		query.Set("after", strconv.FormatInt(*afterMS, 10))
	}
	endpoint.RawQuery = query.Encode()

	body, statusCode, err := httpapi.DoRequest(
		ctx,
		s.httpClient(),
		http.MethodGet,
		endpoint.String(),
		nil,
		map[string]string{
			"Authorization": "Bearer " + accessToken,
			"Accept":        "application/json",
		},
	)
	if err != nil {
		return spotifyRecentlyPlayedResponse{}, err
	}
	if statusCode == http.StatusUnauthorized {
		return spotifyRecentlyPlayedResponse{}, errSpotifyUnauthorized
	}
	if statusCode < 200 || statusCode >= 300 {
		return spotifyRecentlyPlayedResponse{}, fmt.Errorf("spotify recently played failed: status=%d body=%s", statusCode, string(body))
	}

	var response spotifyRecentlyPlayedResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return spotifyRecentlyPlayedResponse{}, err
	}

	return response, nil
}

func (s *Store) fetchSpotifyArtistImages(ctx context.Context, apiURL string, accessToken string, recentlyPlayed spotifyRecentlyPlayedResponse) map[string]*string {
	artistIDs := uniquePrimarySpotifyArtistIDs(recentlyPlayed)
	images := make(map[string]*string, len(artistIDs))

	for _, artistID := range artistIDs {
		artist, err := s.fetchSpotifyArtist(ctx, apiURL, accessToken, artistID)
		if err != nil {
			if s.Logger != nil {
				s.Logger.Printf("Spotify artist image fetch skipped for %s: %v", artistID, err)
			}
			continue
		}

		imageURL := bestSpotifyImageURL(artist.Images)
		if imageURL != "" {
			images[artistID] = &imageURL
		}
	}

	return images
}

func (s *Store) fetchSpotifyArtist(ctx context.Context, apiURL string, accessToken string, artistID string) (spotifyArtist, error) {
	endpoint, err := url.Parse(strings.TrimRight(apiURL, "/") + "/artists/" + url.PathEscape(artistID))
	if err != nil {
		return spotifyArtist{}, err
	}

	body, statusCode, err := httpapi.DoRequest(
		ctx,
		s.httpClient(),
		http.MethodGet,
		endpoint.String(),
		nil,
		map[string]string{
			"Authorization": "Bearer " + accessToken,
			"Accept":        "application/json",
		},
	)
	if err != nil {
		return spotifyArtist{}, err
	}
	if statusCode < 200 || statusCode >= 300 {
		return spotifyArtist{}, fmt.Errorf("spotify artist failed: status=%d body=%s", statusCode, string(body))
	}

	var artist spotifyArtist
	if err := json.Unmarshal(body, &artist); err != nil {
		return spotifyArtist{}, err
	}
	return artist, nil
}

func spotifyPlayHistoryToScrobbleInput(item spotifyPlayHistory, rawPayload []byte, artistImages map[string]*string) (model.ScrobbleInput, bool, error) {
	if item.PlayedAt.IsZero() || strings.TrimSpace(item.Track.Name) == "" || len(item.Track.Artists) == 0 {
		return model.ScrobbleInput{}, false, nil
	}

	primaryArtist := item.Track.Artists[0]
	if strings.TrimSpace(primaryArtist.Name) == "" {
		return model.ScrobbleInput{}, false, nil
	}

	albumImageURL := bestSpotifyImageURL(item.Track.Album.Images)
	req := model.ScrobbleRequest{
		ScrobbledAt:     item.PlayedAt.UTC().Format(time.RFC3339Nano),
		Artist:          primaryArtist.Name,
		ArtistSpotifyID: primaryArtist.ID,
		ArtistURL:       primaryArtist.ExternalURLs.Spotify,
		Album:           item.Track.Album.Name,
		AlbumSpotifyID:  item.Track.Album.ID,
		AlbumURL:        item.Track.Album.ExternalURLs.Spotify,
		AlbumImageURL:   albumImageURL,
		Track:           item.Track.Name,
		TrackSpotifyID:  item.Track.ID,
		TrackURL:        item.Track.ExternalURLs.Spotify,
		TrackImageURL:   albumImageURL,
		Source:          spotifyPollSource,
	}

	if artistImages != nil && primaryArtist.ID != "" && artistImages[primaryArtist.ID] != nil {
		req.ArtistImageURL = *artistImages[primaryArtist.ID]
	}

	input, err := req.ToInput(rawPayload)
	if err != nil {
		return model.ScrobbleInput{}, false, err
	}

	return input, true, nil
}

func uniquePrimarySpotifyArtistIDs(recentlyPlayed spotifyRecentlyPlayedResponse) []string {
	seen := map[string]struct{}{}
	artistIDs := []string{}

	for _, rawItem := range recentlyPlayed.Items {
		var item spotifyPlayHistory
		if err := json.Unmarshal(rawItem, &item); err != nil || len(item.Track.Artists) == 0 {
			continue
		}

		artistID := strings.TrimSpace(item.Track.Artists[0].ID)
		if artistID == "" {
			continue
		}
		if _, ok := seen[artistID]; ok {
			continue
		}

		seen[artistID] = struct{}{}
		artistIDs = append(artistIDs, artistID)
	}

	return artistIDs
}

func bestSpotifyImageURL(images []spotifyImage) string {
	for _, image := range images {
		imageURL := strings.TrimSpace(image.URL)
		if imageURL != "" {
			return imageURL
		}
	}
	return ""
}

func spotifyAccessTokenExpired(expiresAt *int64) bool {
	if expiresAt == nil || *expiresAt <= 0 {
		return true
	}

	const unixSecondsCutoff = int64(100000000000)
	if *expiresAt < unixSecondsCutoff {
		return time.Now().Add(30*time.Second).Unix() >= *expiresAt
	}

	return time.Now().Add(30*time.Second).UnixMilli() >= *expiresAt
}

func (s *Store) httpClient() *http.Client {
	if s.HTTPClient != nil {
		return s.HTTPClient
	}
	return http.DefaultClient
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
