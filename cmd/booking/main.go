package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/infrastructure/env"
	"final-project/intern/booking/infrastructure/logger"
	"final-project/intern/booking/infrastructure/repository"
	"final-project/intern/booking/infrastructure/server"
	"final-project/pkg/booking/constants"

	"log/slog"
)

func parseLogLevel(lvl string) slog.Level {
	switch lvl {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func main() {
	envPath := flag.String("env", "env/booking.dev.env", "path to env file")
	flag.Parse()

	cfg, err := env.NewConfigFromFile(*envPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	lvl := parseLogLevel(cfg.LogLevel)
	lg := logger.New(lvl)
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		lg.Error(constants.MsgFailedOpenDB, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBPing, constants.KeyError, err)
		os.Exit(1)
	}

	client := &http.Client{Timeout: cfg.Hotel.Timeout}
	payClient := &http.Client{Timeout: cfg.Payment.Timeout}
	pg := repository.NewPostgresAdapter(db, lg)
	ht := repository.NewHotelAdapter(cfg.Hotel.BaseURL, client, lg)
	pay := repository.NewPaymentAdapter(cfg.Payment.BaseURL, payClient, lg)
	cr := repository.NewComplexRepo(pg, ht, pay, lg)
	var repo interfaces.Repository = cr

	webhookURL := cfg.Server.Host + cfg.Server.Port + "/webhook/payment"

	lg.Info(constants.EventServerStarted, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventServerStarted, constants.KeyAddr, cfg.Server.Port)
	server.StartServer(cfg.Server.Port, repo, lg, webhookURL)
}
