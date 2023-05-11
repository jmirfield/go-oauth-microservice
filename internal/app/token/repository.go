package token

import (
	"context"
	"oauth/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, token *models.Token) error
	GetByToken(ctx context.Context, token string) (*models.Token, error)
	// DeleteByToken(ctx context.Context, token string) error
}

type tokenRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) (*tokenRepository, error) {
	repo := &tokenRepository{pool}
	err := repo.initTable()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (tr *tokenRepository) initTable() error {
	_, err := tr.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS token (
	id			SERIAL		PRIMARY KEY,
	access		TEXT		NOT NULL,
	expires_at	TIMESTAMPTZ NOT NULL,
	created_at 	TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_token_expires_at ON token (expires_at);
	`)
	return err
}

func (tr *tokenRepository) Create(ctx context.Context, token *models.Token) error {
	_, err := tr.pool.Exec(ctx, `
	INSERT INTO token (access, expires_at)
	VALUES ($1, $2)
	`, token.Access, token.ExpiresAt)
	return err
}

func (tr *tokenRepository) GetByToken(ctx context.Context, token string) (*models.Token, error) {
	var t models.Token

	rows := tr.pool.QueryRow(ctx, "SELECT access, expires_at FROM public.token where access = $1", token)
	err := rows.Scan(
		&t.Access,
		&t.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
