package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"final-project/intern/notification/application"
	"final-project/intern/notification/config"
	"final-project/intern/notification/infrastructure/email"
	"final-project/intern/notification/infrastructure/handlers"
	"final-project/intern/notification/infrastructure/kafka"
	"final-project/intern/notification/infrastructure/metrics"

	"github.com/joho/godotenv"
)

func parseLogLevel(level string) slog.Level {
	switch level {
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
	env := flag.String("env", "", "What env file use")
	levelFlag := flag.String("level", "", "Logging level")
	flag.Parse()

	if *env != "" {
		if err := godotenv.Load(*env); err != nil {
			log.Printf("failed to load env file: %v", err)
		}
	}

	cfg := &config.Config{}
	cfg.LoadFromEnv()

	if *levelFlag != "" {
		cfg.LogLevel = *levelFlag
	}

	logLevel := parseLogLevel(cfg.LogLevel)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("start notification service")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	metricsServer := metrics.AddMetricsHandler(cfg)
	defer metrics.Shutdown(context.Background(), metricsServer)

	var sender application.EmailSender
	if cfg.SMTPUsername == "" || cfg.SMTPPassword == "" {
		slog.Warn("smtp is not configured, using log email sender")
		sender = email.NewLogEmailSender()
	} else {
		sender = email.NewSMTPSender(cfg)
	}
	service := application.NewNotificationService(sender)
	handler := handlers.NewEmailNotificationHandler(service)

	server := kafka.NewServer(cfg, handler.Handle)

	if err := server.Serve(ctx); err != nil {
		slog.Error("notification service stopped with error", "error", err.Error())
		os.Exit(1)
	}
}
