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

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info(logs.MsgStartOperation, logs.KeyEvent, logs.EventServiceInit, slog.String(logs.KeyService, "HotelService"))

	cfg := &config.Config{}
	cfg.LoadFromEnv()

	server.AddMetricsHandler(cfg.Host)
	serverAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := setupDB(cfg, logger)
	if err != nil {
		logger.Error(logs.MsgDatabaseConnectFailed, logs.KeyError, err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(logs.MsgDatabaseCloseFailed, logs.KeyError, err)
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

	logger.Info(logs.MsgStartServer, slog.String("addr", serverAddr))
	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(logs.MsgUnexpectedFail, logs.KeyError, err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info(logs.MsgShuttingDown, logs.KeyEvent, logs.EventServerShutdown)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.HTTPServer.Shutdown(shutdownCtx); err != nil {
		logger.Error(logs.MsgForcedShutdown, logs.KeyError, err)
		os.Exit(1)
	}

	logger.Info(logs.MsgShutdownComplete, logs.KeyEvent, logs.EventServerShutdown)
}

func setupDB(cfg *config.Config, log *slog.Logger) (*sql.DB, error) {
	connStr := cfg.GetDBConnectionString()
	log.Info(logs.MsgStartOperation, logs.KeyEvent, logs.EventDBConnect, slog.String(logs.KeyDBName, cfg.DBName))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info(logs.MsgOperationSuccess, logs.KeyEvent, logs.EventDBConnect)
	return db, nil
}
