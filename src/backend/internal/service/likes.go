package service

import (
	"context"
	"main/internal/model"
	"main/internal/storage"
)

type LikesService struct {
	storage storage.StorageInterface
}

func NewLikesService(storage storage.StorageInterface) *LikesService {
	return &LikesService{storage: storage}
}

func (ls *LikesService) GetLikeCount(ctx context.Context) (int, error) {
	return ls.storage.GetCount(ctx, "likes")
}

// returns the updated likes count, the updated like status for this user, the string representing the action taken, and any error encountered
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
		count, err = ls.storage.DecrementCount(ctx, "likes")
		action = "unliked"
		session.HasLiked = false
	} else {
		count, err = ls.storage.IncrementCount(ctx, "likes")
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
