package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"final-project/intern/booking/domain/models/reservation"
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

	q := r.URL.Query()
	email := q.Get("email")
	hotel := q.Get("hotel")
	numberStr := q.Get("number")
	startStr := q.Get("start")
	endStr := q.Get("end")

	var num uint64
	if _, err := fmt.Sscanf(numberStr, "%d", &num); err != nil {
		logger.Warn("invalid number in webhook URL", constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyError, err)
		http.Error(w, constants.MsgUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		logger.Warn("invalid start in webhook URL", constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyError, err)
		http.Error(w, constants.MsgUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		logger.Warn("invalid end in webhook URL", constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyError, err)
		http.Error(w, constants.MsgUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}

	reserv := reservation.Reserve{
		Email:  email,
		Start:  start,
		End:    end,
		Hotel:  hotel,
		Number: num,
	}

	if p.Status == "timeout" {
		if err := h.service.DeleteReservation(reserv); err != nil {
			logger.Error(constants.EventPaymentWebhook, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventPaymentWebhook)
			http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
