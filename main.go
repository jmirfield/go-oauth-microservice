package main

import (
	"log"
	"oauth/config"
	"oauth/internal/server"
)

func main() {
	cfg := config.LoadConfig()
	err := server.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// http.HandleFunc("/register", app.RegisterHandler())
	// http.HandleFunc("/token", app.TokenHandler())
	// http.HandleFunc("/validate", app.ValidateTokenHandler())
	// http.HandleFunc("/public", app.PublicKeyHandler())

	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}
