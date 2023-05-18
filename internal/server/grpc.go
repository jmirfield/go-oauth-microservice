package server

import (
	"context"
	oauth "oauth/api"
)

func (a *app) GetKey(ctx context.Context, req *oauth.KeyRequest) (*oauth.KeyResponse, error) {
	key, err := a.m.GetPublicKey()
	if err != nil {
		return nil, err
	}
	return &oauth.KeyResponse{Key: string(key)}, nil
}

func (a *app) ValidateToken(ctx context.Context, token *oauth.TokenRequest) (*oauth.TokenResponse, error) {
	return &oauth.TokenResponse{Valid: a.m.ValidateToken(ctx, token.Token)}, nil
}
