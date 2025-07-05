package service

import (
	"context"
	"main/internal/model"
	"main/internal/storage"
)

type VisitorService struct {
	storage storage.StorageInterface
}

func NewVisitorService(storage storage.StorageInterface) *VisitorService {
	return &VisitorService{storage: storage}
}

func (cs *VisitorService) GetVisitorCount(ctx context.Context) (int, error) {
	return cs.storage.GetCount(ctx, "visitors")
}

// IncrementVisitorCount increments the visitor count if the user hasn't visited before
// Returns the updated count and a status message indicating if the count was incremented
func (cs *VisitorService) IncrementVisitorCount(ctx context.Context, session *model.UserSession) (int, string, error) {
	// Check if user has already visited
	if session.HasVisited {
		// User has already visited, just return current count
		count, err := cs.storage.GetCount(ctx, "visitors")
		return count, "already_visited", err
	}

	// User hasn't visited before, increment count
	count, err := cs.storage.IncrementCount(ctx, "visitors")
	if err != nil {
		return 0, "", err
	}

	return count, "incremented", nil
}
