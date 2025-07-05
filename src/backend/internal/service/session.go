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

// GetOrCreateSession returns existing session or creates a new one with default values
func (ss *SessionService) GetOrCreateSession(ctx context.Context, sessionID string) (*model.UserSession, bool, error) {
	if sessionID == "" {
		// Generate new session ID
		sessionID = ss.generateSessionID()
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

	// Try to get existing session
	session, err := ss.storage.GetUserSession(ctx, sessionID)
	if err != nil {
		return nil, false, err
	}

	if session == nil {
		// Session doesn't exist, create new one with existing sessionID
		session = &model.UserSession{
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

	return ss.storage.GetUserSession(ctx, sessionID)
}

func (ss *SessionService) generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
