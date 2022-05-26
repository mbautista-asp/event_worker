package config

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var (
	DBPool    *sqlx.DB
	logFields log.Fields
	Logger    *log.Logger
	dbOnce    sync.Once
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	setUpLogger()
	setUpDb()
}

func setUpDb() {
	dbOnce.Do(
		func() {
			logFields = log.Fields{
				"package":  "config",
				"function": "setUpDb",
			}
			var err error
			dbUrl := os.Getenv("DATABASE_URL")
			connStr := fmt.Sprintf("%v", dbUrl)
			DBPool, err = sqlx.Open("postgres", connStr)
			if err != nil {
				log.WithFields(logFields).Error(err)
			}

			DBPool.SetConnMaxLifetime(0)
			DBPool.SetMaxIdleConns(3)
			DBPool.SetMaxOpenConns(3)

			ctx, stop := context.WithCancel(context.Background())
			defer stop()

			appSignal := make(chan os.Signal, 3)
			signal.Notify(appSignal, os.Interrupt)

			go func() {
				select {
				case <-appSignal:
					stop()
				}
			}()

			ping(ctx)
		},
	)
}

func ping(ctx context.Context) {
	logFields = log.Fields{
		"package":  "config",
		"function": "ping",
	}
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := DBPool.PingContext(ctx); err != nil {
		log.WithFields(logFields).Fatal(err)
	}
}

func setUpLogger() {
	Logger = log.New()
	Logger.SetFormatter(&log.JSONFormatter{})
	Logger.SetReportCaller(true)
	setLogLevel()
}

func setLogLevel() {
	switch level := os.Getenv("LOG_LEVEL"); level {
	case "trace":
		Logger.SetLevel(log.TraceLevel)
	case "debug":
		Logger.SetLevel(log.DebugLevel)
	case "info":
		Logger.SetLevel(log.InfoLevel)
	case "warn":
		Logger.SetLevel(log.WarnLevel)
	case "error":
		Logger.SetLevel(log.ErrorLevel)
	case "fatal":
		Logger.SetLevel(log.FatalLevel)
	case "panic":
		Logger.SetLevel(log.PanicLevel)
	default:
		Logger.SetLevel(log.DebugLevel)
	}
}
