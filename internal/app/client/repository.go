package client

import (
	"context"
	"oauth/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, client *models.Client) error
	GetByID(ctx context.Context, id string) (*models.Client, error)
}

type clientRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) (*clientRepository, error) {
	repo := &clientRepository{pool}

	err := repo.initTable()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (cr *clientRepository) initTable() error {
	_, err := cr.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS client (
	id     TEXT  PRIMARY KEY,
	secret TEXT  NOT NULL
	);
	`)
	return err
}

func (cr *clientRepository) Create(ctx context.Context, client *models.Client) error {
	_, err := cr.pool.Exec(ctx, `
	INSERT INTO client (id, secret)
	VALUES ($1, $2);
	`, client.ID, client.Secret)
	return err
}

func (cr *clientRepository) GetByID(ctx context.Context, id string) (*models.Client, error) {
	var client models.Client

	rows := cr.pool.QueryRow(ctx, "SELECT * FROM public.client WHERE id = $1", id)
	err := rows.Scan(
		&client.ID,
		&client.Secret,
	)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
