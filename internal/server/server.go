package server

import (
	"context"
	"fmt"
	"net/http"
	oauth "oauth/api"
	"oauth/config"
	"oauth/internal/app/client"
	"oauth/internal/app/manager"
	"oauth/internal/app/token"
	"oauth/pkg/rsa"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type app struct {
	m *manager.Manager
	oauth.UnimplementedAuthServer
}

// rootHandler returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from https://github.com/alvarowolfx/go-grpc-rest-api-demo.
func rootHandler(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func Run(cfg *config.Config) error {
	key, err := rsa.GetPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("unable to setup rsa key pairs: %s", err)
	}

	// eventually switch out dbpool for an adapter
	dbpool, err := pgxpool.Connect(context.Background(), cfg.DSN)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %s", err)
	}
	defer dbpool.Close()

	clientRepo, err := client.NewRepository(dbpool)
	if err != nil {
		return fmt.Errorf("failed to setup client repo: %s", err)
	}
	clientService := client.NewService(clientRepo)

	tokenRepo, err := token.NewRepository(dbpool)
	if err != nil {
		return fmt.Errorf("failed to setup token repo: %s", err)
	}
	defer tokenRepo.Close()
	tokenService := token.NewService(tokenRepo, key)

	manager := manager.NewManager(clientService, tokenService)

	app := &app{m: manager}

	r := chi.NewRouter()
	app.setupRoutes(r, "v1")

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	oauth.RegisterAuthServer(grpcServer, app)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: rootHandler(grpcServer, r),
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
		}
	}()
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *app) setupRoutes(r chi.Router, version string) {
	r.Route(fmt.Sprintf("/%s", version), func(r chi.Router) {
		r.Post("/register", a.registerHandler())
		r.Get("/token", a.tokenHandler())
		r.Get("/validate", a.tokenValidationHandler())
	})
}
