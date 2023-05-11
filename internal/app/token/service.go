package token

import (
	"context"
	"crypto/rsa"
	"oauth/internal/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type Service interface {
	Create(ctx context.Context, client *models.Client) (*models.Token, error)
	GetAccess(ctx context.Context, token string) (*models.Token, error)
	Public() *rsa.PublicKey
}

type tokenService struct {
	r Repository
	k *rsa.PrivateKey
}

func NewService(repo Repository, key *rsa.PrivateKey) *tokenService {
	return &tokenService{repo, key}
}

func (ts *tokenService) Create(ctx context.Context, client *models.Client) (*models.Token, error) {
	exp := time.Now().Add(10 * time.Minute)
	claims := jwt.StandardClaims{
		Audience:  client.ID,
		ExpiresAt: exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	access, err := token.SignedString(ts.k)
	if err != nil {
		return nil, err
	}

	t := &models.Token{Access: access, ExpiresAt: exp}
	err = ts.r.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (ts *tokenService) GetAccess(ctx context.Context, token string) (*models.Token, error) {
	return ts.r.GetByToken(ctx, token)
}

func (ts *tokenService) Public() *rsa.PublicKey {
	return &ts.k.PublicKey
}
