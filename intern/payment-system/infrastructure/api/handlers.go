package api

import (
	"encoding/json"
	"final-project/intern/payment-system/application"
	"final-project/intern/payment-system/constants"
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

	slog.Info(constants.MsgCreateLink,
		constants.KeyWebHook, requestParams.WebHook,
		constants.KeyAmount, requestParams.Amount,
		constants.KeyMessage, requestParams.Message,
	)

	ctx := request.Context()

	url, err := manager.UseCase.CreatePaymentLink(ctx, requestParams.Amount, requestParams.WebHook, requestParams.Message)
	if err != nil {
		slog.Error(constants.MsgFailedLinkCreation, constants.KeyError, err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(
		LinkResponse{
			Link: url,
		},
	)
	if err != nil {
		slog.Error(constants.MsgFailedToCreateResponse, constants.KeyError, err.Error())
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

	slog.Info(constants.MsgRenderPage,
		constants.KeyPaymentKey, key,
	)

	err := manager.UseCase.RenderPaymentPage(request.Context(), key, writer)

	if err != nil {
		slog.Error(constants.MsgFailedRenderPage, constants.KeyError, err.Error(), constants.KeyPaymentKey, key)
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

	slog.Info(constants.MsgConfirm,
		constants.KeyPaymentKey, key,
	)

	err := manager.UseCase.ConfirmPayment(request.Context(), key)

	if err != nil {
		slog.Error(constants.MsgFailedConfirm, constants.KeyError, err.Error(), constants.KeyPaymentKey, key)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}
