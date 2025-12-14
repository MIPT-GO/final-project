package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"final-project/intern/hotel/infrastructure/server"
	config "final-project/pkg/hotel"
	"final-project/pkg/logs"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info(logs.MsgStartOperation, logs.KeyEvent, logs.EventServiceInit, slog.String(logs.KeyService, "HotelService"))

	cfg := &config.Config{}
	cfg.LoadFromEnv()

	serverAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := setupDB(cfg, logger)
	if err != nil {
		logger.Error("Failed to connect to database", logs.KeyError, err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close DB connection", logs.KeyError, err)
		}
	}()

	srv := server.NewServer(
		logger,
		db,
		serverAddr,
		time.Duration(cfg.ReadTimeout)*time.Second,
		time.Duration(cfg.WriteTimeout)*time.Second,
		time.Duration(cfg.IdleTimeout)*time.Second,
	)

	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(logs.MsgUnexpectedFail, logs.KeyError, err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.HTTPServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", logs.KeyError, err)
		os.Exit(1)
	}

}

func setupDB(cfg *config.Config, log *slog.Logger) (*sql.DB, error) {
	connStr := cfg.GetDBConnectionString()
	log.Info(logs.MsgStartOperation, logs.KeyEvent, logs.EventDBQuery, slog.String("db_name", cfg.DBName))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
