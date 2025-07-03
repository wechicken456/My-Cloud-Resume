package service

import (
	"context"
	"main/internal/model"
	"main/internal/storage"
)

type LikesService struct {
	storage *storage.Storage
}

func NewLikesService(storage *storage.Storage) *LikesService {
	return &LikesService{storage: storage}
}

func (ls *LikesService) GetLikeCount(ctx context.Context) (int, error) {
	return ls.storage.GetCount("likes")
}

// Increments/Decrements the like count and toggles the user's like status
func (ls *LikesService) ToggleLike(ctx context.Context, sessionID string) (int, bool, string, error) {
	// Check current like status
	session, err := ls.storage.GetUserSession(ctx, sessionID)
	if err != nil {
		return 0, false, "", err
	}

	if session == nil {
		// Create new session
		err = ls.storage.CreateUserSession(ctx, sessionID)
		if err != nil {
			return 0, false, "", err
		}
		session = &model.UserSession{SessionID: sessionID, HasLiked: false}
	}

	var count int
	var action string

	if session.HasLiked {
		count, err = ls.storage.DecrementCount("likes")
		action = "unliked"
		session.HasLiked = false
	} else {
		count, err = ls.storage.IncrementCount("likes")
		action = "liked"
		session.HasLiked = true
	}

	if err != nil {
		return 0, false, "", err
	}

	err = ls.storage.UpdateUserSession(ctx, session)
	if err != nil {
		return count, session.HasLiked, action, err
	}

	return count, session.HasLiked, action, nil
}
