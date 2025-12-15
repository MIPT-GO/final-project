package api

import (
	"encoding/json"
	"final-project/intern/payment-system/application"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type HandlersManager struct {
	UseCase *application.PaymentsUseCase
}

func (manager *HandlersManager) PaymentRequestHandler(writer http.ResponseWriter, request *http.Request) {
	var requestParams PaymentRequest
	err := json.NewDecoder(request.Body).Decode(&requestParams)
	defer request.Body.Close()

	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	slog.Info("Create payment link request",
		"webhook", requestParams.WebHook,
		"amount", requestParams.Amount,
		"message", requestParams.Message,
	)

	ctx := request.Context()

	url, err := manager.UseCase.CreatePaymentLink(ctx, requestParams.Amount, requestParams.WebHook, requestParams.Message)
	if err != nil {
		slog.Error("Error during create payment link", "error", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(
		LinkResponse{
			Link: url,
		},
	)
	if err != nil {
		slog.Error("Error during creating response", "error", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (manager *HandlersManager) RenderPaymentPageHandler(writer http.ResponseWriter, request *http.Request) {
	key := mux.Vars(request)["key"]

	if key == "" {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	slog.Info("Render payment page request",
		"key", key,
	)

	err := manager.UseCase.RenderPaymentPage(request.Context(), key, writer)

	if err != nil {
		slog.Error("Error during render payment page", "error", err.Error(), "key", key)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (manager *HandlersManager) ConfirmPaymentHandler(writer http.ResponseWriter, request *http.Request) {
	key := mux.Vars(request)["key"]

	if key == "" {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	slog.Info("Confirm payment request",
		"key", key,
	)

	err := manager.UseCase.ConfirmPayment(request.Context(), key)

	if err != nil {
		slog.Error("Error during confirm payment", "error", err.Error(), "key", key)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}
