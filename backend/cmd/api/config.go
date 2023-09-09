package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func updateConfigWithVariables() (*config, error) {
	var cfg config

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// SERVER Config
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	flag.IntVar(&cfg.port, "port", port, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// DATABASE Config
	maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	maxOpenConns, err := strconv.Atoi(maxOpenConnsStr)
	if err != nil {
		log.Fatal(err)
	}
	maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS")
	maxIdleConns, err := strconv.Atoi(maxIdleConnsStr)
	if err != nil {
		log.Fatal(err)
	}

	// DB Config
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("CFBOX_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", maxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", maxIdleConns, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", os.Getenv("DB_MAX_IDLE_TIME"), "PostgreSQL max connection idle time")

	// Email Config
	emailPortStr := os.Getenv("EMAIL_SERVER_PORT")
	emailPort, err := strconv.Atoi(emailPortStr)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("EMAIL_HOST_SERVER"), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", emailPort, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("EMAIL_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("EMAIL_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("EMAIL_SENDER"), "SMTP sender")

	// Redis config
	flag.StringVar(&cfg.redisURL, "redis-url", os.Getenv("REDIS_URL"), "Redis URL")

	// Token Expiration
	tokenExpirationStr := os.Getenv("TOKEN_EXPIRATION")
	duration, err := time.ParseDuration(tokenExpirationStr)
	if err != nil {
		return nil, err
	}
	cfg.tokenExpiration.durationString = tokenExpirationStr
	cfg.tokenExpiration.duration = duration

	// Frontend URL
	flag.StringVar(&cfg.frontendURL, "frontend-url", os.Getenv("FRONTEND_URL"), "Frontend URL")

	// CORS
	// flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(s string) error {
	// 	cfg.cors.trustedOrigins = strings.Fields(s)
	// 	return nil
	// })

	cfg.cors.trustedOrigins = []string{"http://localhost:5173"}

	flag.Parse()

	return &cfg, nil
}
