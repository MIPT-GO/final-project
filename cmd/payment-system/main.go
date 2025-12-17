package main

import (
	"context"
	"final-project/intern/payment-system/application"
	"final-project/intern/payment-system/config"
	"final-project/intern/payment-system/constants"
	"final-project/intern/payment-system/infrastructure/api"
	"final-project/intern/payment-system/infrastructure/render"
	"final-project/intern/payment-system/infrastructure/storage"
	"final-project/intern/payment-system/infrastructure/webhook"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	env := flag.String("env", "", "What env file use")
	log_level := flag.String("level", "info", "Logging level")

	flag.Parse()

	var level slog.Level
	switch *log_level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)

	if *env != "" {
		godotenv.Load(*env)
	}

	config := config.Config{}
	config.LoadFromEnv()

	redis := storage.RedisStorage{}
	err := redis.NewStorage(context.Background(), &config)

	if err != nil {
		logger.Error(constants.MsgFailedConnectionDB, constants.KeyError, err.Error())
		panic(err)
	}

	generator := storage.RandomKeyGenerator{}

	render := render.HTMLRender{}
	render.New(&config)

	sender := webhook.NetWebHookSender{
		Client: &http.Client{},
	}

	usecase := application.PaymentsUseCase{
		Repository: &redis,
		Generator:  &generator,
		Render:     &render,
		Sender:     &sender,
		Config:     &config,
	}

	api.StartServer(context.Background(), &usecase, &config)
}
