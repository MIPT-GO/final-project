package server

import (
	"context"
	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/application/service"
	"final-project/intern/booking/infrastructure/server/handlers"
	"final-project/pkg/booking/constants"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func StartServer(port string, repo interfaces.Repository, logger interfaces.Logger, webhook string, producer interfaces.Producer) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s := service.NewReservationService(repo, webhook, logger, producer)
	h := handlers.NewHandler(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/booking/email/{email}/", h.GetByEmail)
	mux.HandleFunc("/v1/booking/hotel/{hotel}/", h.GetByHotel)
	mux.HandleFunc("/v1/booking/", h.BookRoomInHotel)
	mux.HandleFunc("/v1/available/{hotel}/", h.Available)
	mux.HandleFunc("/webhook/payment/{id}", h.PaymentWebhook)
	serv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		logger.Info(constants.EventServerStarted, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventServerStarted, constants.KeyAddr, serv.Addr)
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(constants.EventServerStarted, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventServerStarted)
		}
	}()

	<-ctx.Done()
	logger.Info(constants.EventServerShutdown, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventServerShutdown)
	if err := serv.Shutdown(context.Background()); err != nil {
		logger.Error(constants.EventServerShutdown, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventServerShutdown)
	}
}
