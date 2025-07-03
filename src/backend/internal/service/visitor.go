package service

import (
	"context"
	"main/internal/model"
	"main/internal/storage"
)

type CounterService struct {
	storage *storage.Storage
}

func NewCounterService(storage *storage.Storage) *CounterService {
	return &CounterService{storage: storage}
}

func (cs *CounterService) GetVisitorCount(ctx context.Context) (int, error) {
	return cs.storage.GetCount("visitor")
}

func (cs *CounterService) IncrementVisitorCount(ctx context.Context, sessionID string) (int, bool, string, error) {
	// Check if user has already visited
	session, err := cs.storage.GetUserSession(ctx, sessionID)
	if err != nil {
		return 0, false, "", err
	}

	if session == nil {
		// Create new session
		err = cs.storage.CreateUserSession(ctx, sessionID)
		if err != nil {
			return 0, false, "", err
		}
		session = &model.UserSession{SessionID: sessionID, HasVisited: false}
	}

	if session.HasVisited {
		// User already visited, just return current count
		count, err := cs.storage.GetCount("visitor")
		return count, false, "already_counted", err
	}

	// Increment count and mark as visited
	count, err := cs.storage.IncrementCount("visitor")
	if err != nil {
		return 0, false, "", err
	}

	session.HasVisited = true
	err = cs.storage.UpdateUserSession(ctx, session)
	if err != nil {
		return count, true, "", err // Return count even if session update fails
	}

	return count, true, "incremented", nil
}

func (cs *CounterService) GetSessionStatus(ctx context.Context, sessionID string) (*model.UserSession, error) {
	session, err := cs.storage.GetUserSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		// Return default session
		return &model.UserSession{
			SessionID:  sessionID,
			HasVisited: false,
			HasLiked:   false,
		}, nil
	}

	return session, nil
}
