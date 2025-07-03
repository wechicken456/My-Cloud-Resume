package service

import (
	"context"
	"main/internal/config"
	"main/internal/model"
)

type ContactService struct {
	config *config.Config
}

func NewContactService(cfg *config.Config) *ContactService {
	return &ContactService{config: cfg}
}

func (cs *ContactService) ProcessContactRequest(ctx context.Context, contactReq *model.ContactRequest) error {
	return nil
}
