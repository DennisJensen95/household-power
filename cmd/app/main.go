package main

import (
	"log"

	"github.com/DennisJensen95/golang-rest-api/config"
	"github.com/DennisJensen95/golang-rest-api/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
