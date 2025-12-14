package handlers

import (
	"encoding/json"
	"final-project/intern/booking/infrastructure/server/dto"
	"final-project/pkg/booking/constants"
	"net/http"
)

func (handl *handler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	logger := handl.service.Logger

	logger.Info(
		constants.EventIncomingRequest,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventIncomingRequest,
		constants.KeyHandler, "GetByEmail",
		constants.KeyMethod, r.Method,
		constants.KeyPath, r.URL.Path,
	)

	if r.Method != http.MethodGet {
		logger.Warn(
			constants.EventMethodNotAllowed,
			constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventMethodNotAllowed,
			constants.KeyHandler, "GetByEmail",
			constants.KeyMethod, r.Method,
		)

		http.Error(w, constants.MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	email := r.PathValue("email")

	logger.Debug(
		constants.EventGetByEmail,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventGetByEmail,
		constants.KeyUserEmail, email,
	)

	reservs, err := handl.service.GetByEmail(email)
	if err != nil {
		logger.Error(
			constants.EventFailedToFetch,
			err,
			constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch,
			constants.KeyUserEmail, email,
		)
		http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		return
	}

	logger.Info(
		constants.EventRequestCompleted,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventRequestCompleted,
		constants.KeyUserEmail, email,
		constants.KeyCount, len(reservs),
	)

	resp := make([]dto.Reservation, 0, len(reservs))
	for _, reserv := range reservs {
		resp = append(resp, dto.Reservation{
			Email:  reserv.Email,
			Start:  reserv.Start.String(),
			End:    reserv.End.String(),
			Hotel:  reserv.Hotel,
			Number: reserv.Number,
		})
	}

	out := dto.ReservationsResponse{
		Reservations: resp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(out); err != nil {
		logger.Error(
			constants.EventFailedToEncode,
			err,
			constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToEncode,
			constants.KeyUserEmail, email,
		)
		http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		return
	}

	logger.Info(
		constants.EventRequestCompleted,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventRequestCompleted,
		constants.KeyHandler, "GetByEmail",
		constants.KeyUserEmail, email,
		constants.KeyStatus, http.StatusOK,
	)
}
