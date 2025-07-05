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

// ToggleLike toggles the like status for a session and returns the updated count and action taken
func (ls *LikesService) ToggleLike(ctx context.Context, session *model.UserSession) (int, string, error) {
	var count int
	var action string
	var err error

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
		return 0, "", err
	}

	return count, action, nil
}
