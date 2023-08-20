package main

import "flag"

func updateConfigWithVariables() (*config, error) {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://cfbox:pa55w0rd@localhost/cfbox?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	return &cfg, nil
}
