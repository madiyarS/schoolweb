package services

import (
	"github.com/gorilla/sessions"
)

type SessionService struct {
	store *sessions.CookieStore
}

func NewSessionService(key []byte) *SessionService {
	store := sessions.NewCookieStore(key)
	return &SessionService{store: store}
}

func (s *SessionService) GetStore() *sessions.CookieStore {
	return s.store
}
