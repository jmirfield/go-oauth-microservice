package manager

import (
	"context"
	"fmt"
	"oauth/internal/app/client"
	"oauth/internal/app/token"
	"oauth/internal/errors"
	"oauth/internal/models"
)

// Manager orchestrates client and token services
type Manager struct {
	clientService client.Service
	tokenService  token.Service
}

// NewManager -
func NewManager(cs client.Service, ts token.Service) *Manager {
	return &Manager{cs, ts}
}

// RegisterClient handles client registration
func (m *Manager) RegisterClient(ctx context.Context) (*models.Client, error) {
	client, err := m.clientService.Create(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, errors.ErrInternalServer
	}

	return client, nil
}

// GenerateToken handles token generation
func (m *Manager) GenerateToken(ctx context.Context, req *models.Client) (*models.Token, error) {
	client, err := m.clientService.GetByID(ctx, req.ID)
	if err != nil || req.Secret != client.Secret {
		return nil, errors.ErrInvalidClient
	}

	token, err := m.tokenService.Create(ctx, client)
	if err != nil {
		fmt.Println(err)
		return nil, errors.ErrInternalServer
	}

	return token, nil
}
