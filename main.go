package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"oauth/config"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	cfg := config.LoadConfig()

	key, err := os.ReadFile(cfg.Key)
	if err != nil {
		log.Fatalf("Unable to read key file: %v\n", err)
	}

	dbpool, err := pgxpool.Connect(context.Background(), cfg.DSN)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	app, err := NewApp(key, dbpool)
	if err != nil {
		log.Fatalf("Unable to setup app: %s\n", err)
	}
	http.HandleFunc("/register", app.RegisterHandler())
	http.HandleFunc("/token", app.TokenHandler())
	http.HandleFunc("/validate", app.ValidateTokenHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}
