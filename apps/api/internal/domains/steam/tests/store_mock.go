package tests

import "context"

type storeMock struct {
	message      *string
	err          error
	receivedID   string
	receivedName string
}

func (s *storeMock) GetPlayerCount(_ context.Context, gameID string, gameName string) (*string, error) {
	s.receivedID = gameID
	s.receivedName = gameName
	return s.message, s.err
}
