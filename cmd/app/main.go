package main

import (
	"log"

	"github.com/andreyxaxa/Delayed-Notifier/config"
	"github.com/andreyxaxa/Delayed-Notifier/internal/app"
	"github.com/joho/godotenv"
)

// TODO: docker-compose - ДЕЛАТЬ HEALTHCHECK RABBITMQ

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("config error: %s", err)
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}
