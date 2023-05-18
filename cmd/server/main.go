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
}
