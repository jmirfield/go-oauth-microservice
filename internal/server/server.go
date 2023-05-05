package server

import (
	"context"
	"fmt"
	"net/http"
	"oauth/config"
	"oauth/internal/app/client"
	"oauth/internal/app/manager"
	"oauth/internal/app/token"
	"oauth/pkg/rsa"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

type server struct {
	m *manager.Manager
}

func Run(cfg *config.Config) error {
	key, err := rsa.GetPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("unable to setup rsa key pairs: %s", err)
	}

	dbpool, err := pgxpool.Connect(context.Background(), cfg.DSN)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %s", err)
	}
	defer dbpool.Close()

	// eventually switch out dbpool for an adapter
	clientRepo, err := client.NewRepository(dbpool)
	if err != nil {
		return fmt.Errorf("failed to setup client repo: %s", err)
	}
	clientService := client.NewService(clientRepo)

	// eventually switch out dbpool for an adapter
	tokenRepo, err := token.NewRepository(dbpool)
	if err != nil {
		return fmt.Errorf("failed to setup token repo: %s", err)
	}
	tokenService := token.NewService(tokenRepo, key)

	manager := manager.NewManager(clientService, tokenService)

	srv := &server{manager}

	r := chi.NewRouter()
	srv.setupRoutes(r, "v1")

	http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r)
	return nil
}

func (s *server) setupRoutes(r chi.Router, version string) {
	r.Route(fmt.Sprintf("/%s", version), func(r chi.Router) {
		r.Post("/register", s.registerHandler())
		r.Get("/token", s.tokenHandler())
	})
}
