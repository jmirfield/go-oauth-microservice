package client

import (
	"context"
	"oauth/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context) (*models.Client, error)
	GetByID(ctx context.Context, id string) (*models.Client, error)
}

type clientService struct {
	r Repository
}

func NewService(repo Repository) *clientService {
	return &clientService{repo}
}

func (cs *clientService) Create(ctx context.Context) (*models.Client, error) {
	clientID := uuid.New().String()
	clientSecret := uuid.New().String()
	client := &models.Client{
		ID:     clientID,
		Secret: clientSecret,
	}

	err := cs.r.Create(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (cs *clientService) GetByID(ctx context.Context, id string) (*models.Client, error) {
	client, err := cs.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return client, nil
}
