package main

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"time"

	"crossfitbox.booking.system/internal/data"
	"crossfitbox.booking.system/internal/jsonlog"
	"crossfitbox.booking.system/internal/mailer"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	redisURL        string
	tokenExpiration struct {
		durationString string
		duration       time.Duration
	}
	secret struct {
		HMC               string
		secretKey         []byte
		sessionExpiration time.Duration
	}
	frontendURL string
	cors        cors.Options
}

type application struct {
	config      config
	logger      *jsonlog.Logger
	models      data.Models
	mailer      mailer.Mailer
	redisClient *redis.Client
	wg          sync.WaitGroup
}

func main() {

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := updateConfigWithVariables()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	db, err := openDB(*cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	// sdkConfig := aws.Config{
	// 	Region: cfg.awsConfig.Region,
	// 	Credentials: credentials.NewStaticCredentialsProvider(
	// 		cfg.awsConfig.AccessKeyID, cfg.awsConfig.AccessKeySecret, "",
	// 	),
	// }

	redisClient, err := openRedis(*cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	logger.PrintInfo("redis database connection established", nil)

	app := &application{
		config:      *cfg,
		logger:      logger,
		models:      data.NewModels(db),
		mailer:      mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		redisClient: redisClient,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func openRedis(cfg config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ctx := context.Background()

	err = client.Ping(ctx).Err()
	if err != nil {

		return nil, err
	}
	return client, nil
}
