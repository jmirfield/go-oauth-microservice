package main

import (
	"log"
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
func NewApp(key []byte, db *pgxpool.Pool) *App {
	manager := manage.NewManager()
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", key, jwt.SigningMethodRS256))

	adapter := pgx4adapter.NewPool(db)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()

	clientStore, _ := pg.NewClientStore(adapter)

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
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	return &App{
		cs:   clientStore,
		osrv: srv,
	}
}
