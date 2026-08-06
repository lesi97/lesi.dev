package tests

import (
	"context"
	"time"

	auth_handler "github.com/lesi97/lesi.dev/internal/domains/auth/handler"
	auth_store "github.com/lesi97/lesi.dev/internal/domains/auth/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type authStoreStub struct {
	anilistAuthURL                   *string
	anilistAuthErr                   error
	anilistCallbackErr               error
	spotifyAuthURL                   *string
	spotifyAuthErr                   error
	spotifyCallbackErr               error
	twitchModAuthURL                 *string
	twitchModAuthErr                 error
	twitchModCallbackErr             error
	twitchFrontendAuthStartResult    *auth_store.TwitchFrontendAuthStartResult
	twitchFrontendAuthStartErr       error
	twitchFrontendCallbackIdentity   *auth_store.TwitchFrontendIdentity
	twitchFrontendCallbackErr        error
	twitchFrontendUpsertUserID       string
	twitchFrontendUpsertUserErr      error
	twitchFrontendCreateSessionToken string
	twitchFrontendCreateSessionErr   error
	twitchFrontendGetUser            *auth_store.TwitchFrontendUser
	twitchFrontendGetUserErr         error
	twitchFrontendDeleteSessionErr   error
}

func (s *authStoreStub) AnilistAuthUrl() (*string, error) {
	return s.anilistAuthURL, s.anilistAuthErr
}

func (s *authStoreStub) AnilistCallback(ctx context.Context, code string) error {
	return s.anilistCallbackErr
}

func (s *authStoreStub) SpotifyAuthUrl(ctx context.Context) (*string, error) {
	return s.spotifyAuthURL, s.spotifyAuthErr
}

func (s *authStoreStub) SpotifyCallback(ctx context.Context, code string) error {
	return s.spotifyCallbackErr
}

func (s *authStoreStub) TwitchModAuthUrl() (*string, error) {
	return s.twitchModAuthURL, s.twitchModAuthErr
}

func (s *authStoreStub) TwitchModCallback(ctx context.Context, code string) error {
	return s.twitchModCallbackErr
}

func (s *authStoreStub) TwitchFrontendAuthStart() (*auth_store.TwitchFrontendAuthStartResult, error) {
	return s.twitchFrontendAuthStartResult, s.twitchFrontendAuthStartErr
}

func (s *authStoreStub) TwitchFrontendCallback(ctx context.Context, code string, state string, expectedState string, pkceVerifier string) (*auth_store.TwitchFrontendIdentity, error) {
	return s.twitchFrontendCallbackIdentity, s.twitchFrontendCallbackErr
}

func (s *authStoreStub) TwitchFrontendUpsertUser(ctx context.Context, identity auth_store.TwitchFrontendIdentity) (string, error) {
	return s.twitchFrontendUpsertUserID, s.twitchFrontendUpsertUserErr
}

func (s *authStoreStub) TwitchFrontendCreateSession(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	return s.twitchFrontendCreateSessionToken, s.twitchFrontendCreateSessionErr
}

func (s *authStoreStub) TwitchFrontendGetUserBySession(ctx context.Context, sessionToken string) (*auth_store.TwitchFrontendUser, error) {
	return s.twitchFrontendGetUser, s.twitchFrontendGetUserErr
}

func (s *authStoreStub) TwitchFrontendDeleteSessionByToken(ctx context.Context, sessionToken string) error {
	return s.twitchFrontendDeleteSessionErr
}

func newAuthHandler(logger *utils.Logger, store auth_store.Methods) *auth_handler.Handler {
	return auth_handler.NewHandlerWithStore(logger, store)
}
