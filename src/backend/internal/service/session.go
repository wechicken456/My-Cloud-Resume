package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"main/internal/model"
	"main/internal/storage"
)

type SessionService struct {
	storage storage.StorageInterface
}

func NewSessionService(storage storage.StorageInterface) *SessionService {
	return &SessionService{storage: storage}
}

// createNewSession creates a new default session with server-generated ID
func (ss *SessionService) createNewSession(ctx context.Context) (*model.UserSession, bool, error) {
	sessionID := ss.generateSessionID()
	session := &model.UserSession{
		SessionID:  sessionID,
		HasVisited: false,
		HasLiked:   false,
	}

	err := ss.storage.CreateUserSession(ctx, sessionID)
	if err != nil {
		return nil, false, err
	}

	return session, true, nil // true indicates new session was created
}

// GetOrCreateSession returns existing session or creates a new one with default values
func (ss *SessionService) GetOrCreateSession(ctx context.Context, sessionID string) (*model.UserSession, bool, error) {
	if sessionID == "" {
		return ss.createNewSession(ctx) // No existing session, create a new one
	}

	// Try to get existing session
	session, err := ss.storage.GetUserSession(ctx, sessionID)
	if err != nil {
		return nil, false, err
	}

	// If session is nil, create a new one
	if session == nil {
		return ss.createNewSession(ctx) // No existing session, create a new one
	}
	return session, false, nil // false indicates existing session was found
}

// UpdateSession updates an existing session
func (ss *SessionService) UpdateSession(ctx context.Context, session *model.UserSession) error {
	return ss.storage.UpdateUserSession(ctx, session)
}

// ValidateSession checks if a session exists and is valid
func (ss *SessionService) ValidateSession(ctx context.Context, sessionID string) (*model.UserSession, error) {
	if sessionID == "" {
		return nil, nil
	}
	res, err := ss.storage.GetUserSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (ss *SessionService) generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
