package manager

import (
	"context"
	"fmt"
	"oauth/internal/app/client"
	"oauth/internal/app/token"
	"oauth/internal/errors"
	"oauth/internal/models"
	"oauth/pkg/rsa"
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
func (m *Manager) GenerateToken(ctx context.Context, reqClient *models.Client) (*models.Token, error) {
	client, err := m.clientService.GetByID(ctx, reqClient.ID)
	if err != nil || reqClient.Secret != client.Secret {
		return nil, errors.ErrInvalidClient
	}

	token, err := m.tokenService.Create(ctx, client)
	if err != nil {
		fmt.Println(err)
		return nil, errors.ErrInternalServer
	}

	return token, nil
}

func (m *Manager) ValidateToken(ctx context.Context, reqToken string) bool {
	_, err := m.tokenService.GetAccess(ctx, reqToken)
	if err != nil {
		return false
	}

	return true
}

func (m *Manager) GetPublicKey() ([]byte, error) {
	return rsa.PublicBytes(m.tokenService.Public())
}
