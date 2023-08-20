package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func updateConfigWithVariables() (*config, error) {
	var cfg config

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("CFBOX_DB_DSN"), "PostgreSQL DSN")

	flag.Parse()

	return &cfg, nil
}
