package api

import (
	"context"
	"final-project/intern/payment-system/application"
	"final-project/intern/payment-system/config"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateServer(context context.Context, usecase *application.PaymentsUseCase, config *config.Config) *http.Server {
	router := mux.NewRouter()

	handlers := HandlersManager{
		UseCase: usecase,
	}

	router.HandleFunc(config.RouterPrefix+"/link", handlers.PaymentRequestHandler).Methods("POST")
	router.HandleFunc(config.RouterPrefix+"/confirm/{key}", handlers.ConfirmPaymentHandler).Methods("POST")
	router.HandleFunc(config.RouterPrefix+"/page/{key}", handlers.RenderPaymentPageHandler).Methods("GET")

	server := http.Server{
		Handler: router,
		Addr:    config.Host + ":" + strconv.Itoa(config.Port),
	}

	return &server
}

func StartServer(context context.Context, usecase *application.PaymentsUseCase, config *config.Config) {
	server := CreateServer(context, usecase, config)

	ctx, cancel := signal.NotifyContext(context, os.Interrupt)
	defer cancel()

	slog.Info("Start payment system")
	go func() {
		server.ListenAndServe()
	}()
	<-ctx.Done()
	slog.Info("Shutdown payment system")

	err := server.Shutdown(context)
	if err != nil {
		fmt.Println(err.Error())
	}
}
