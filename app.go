package main

import (
	"fmt"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

// App struct
type App struct {
	cs   *pg.ClientStore
	osrv *server.Server
}

// NewApp constructor
func NewApp(key []byte, db *pgxpool.Pool) (*App, error) {
	manager := manage.NewManager()
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", key, jwt.SigningMethodRS256))

	adapter := pgx4adapter.NewPool(db)
	tokenStore, err := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	if err != nil {
		return nil, fmt.Errorf("failed to setup token store: %s", err)
	}
	defer tokenStore.Close()

	clientStore, err := pg.NewClientStore(adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to setup client store: %s", err)
	}

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	srvCfg := server.NewConfig()
	srvCfg.AllowedGrantTypes = []oauth2.GrantType{
		oauth2.ClientCredentials,
	}

	srv := server.NewServer(srvCfg, manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		fmt.Printf("Internal Error: %s\n", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		fmt.Printf("Response Error: %s\n", re.Error.Error())
	})

	return &App{
		key:  &privateKey.PublicKey,
		cs:   clientStore,
		osrv: srv,
	}, nil
}
