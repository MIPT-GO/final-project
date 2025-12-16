package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"final-project/pkg/booking/constants"
)

type paymentWebhookPayload struct {
	Status string `json:"status"`
}

func (h *handler) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	logger := h.service.Logger
	logger.Info(constants.EventPaymentWebhook, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventPaymentWebhook, constants.KeyMethod, r.Method, constants.KeyPath, r.URL.Path)
	var p paymentWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		logger.Warn("failed to decode payment webhook status", constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyError, err)
		http.Error(w, constants.MsgUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}

	// extract id from path: /webhook/payment/{id}
	path := strings.TrimPrefix(r.URL.Path, "/webhook/payment/")
	path = strings.Trim(path, "/")
	if path == "" {
		logger.Warn("missing id in webhook path", constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch)
		http.Error(w, constants.MsgUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}
	id, err := strconv.ParseUint(path, 10, 64)
	if err != nil {
		logger.Warn("invalid id in webhook path", constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyError, err)
		http.Error(w, constants.MsgUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}

	reserv, err := h.service.GetById(uint64(id))
	if err != nil {
		logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch)
		http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		return
	}

	if p.Status == "timeout" {
		if err := h.service.DeleteReservation(reserv); err != nil {
			logger.Error(constants.EventPaymentWebhook, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventPaymentWebhook)
			http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
			return
		}
	}
	if h.service.Producer != nil {
		ownerEmail, _ := h.service.Repo().GetHotelOwnerEmail(reserv.Hotel)
		e := map[string]any{
			"event":       "booking_created",
			"email":       reserv.Email,
			"owner_email": ownerEmail,
			"hotel":       reserv.Hotel,
			"number":      reserv.Number,
			"start":       reserv.Start.Format(time.RFC3339),
			"end":         reserv.End.Format(time.RFC3339),
		}
		b, _ := json.Marshal(e)
		h.service.Producer.Send(b)
		logger.Debug("Start to send message to kafka")
	}
	w.WriteHeader(http.StatusOK)
}
