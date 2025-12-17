package handlers

import (
	"encoding/json"
	"final-project/intern/booking/constants"
	"final-project/intern/booking/infrastructure/server/dto"
	"net/http"
)

func (handl *handler) GetByHotel(w http.ResponseWriter, r *http.Request) {
	logger := handl.service.Logger

	logger.Info(
		constants.EventIncomingRequest,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventIncomingRequest,
		constants.KeyHandler, "GetByHotel",
		constants.KeyMethod, r.Method,
		constants.KeyPath, r.URL.Path,
	)

	if r.Method != http.MethodGet {
		logger.Warn(
			constants.EventMethodNotAllowed,
			constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventMethodNotAllowed,
			constants.KeyHandler, "GetByHotel",
			constants.KeyMethod, r.Method,
		)

		http.Error(w, constants.MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	hotel := r.PathValue("hotel")

	logger.Debug(
		constants.EventGetByHotel,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventGetByHotel,
		constants.KeyHotelName, hotel,
	)

	reservs, err := handl.service.GetByHotel(hotel)
	if err != nil {
		logger.Error(
			constants.EventFailedToFetch,
			err,
			constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch,
			constants.KeyHotelName, hotel,
		)
		http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		return
	}

	logger.Info(
		constants.EventRequestCompleted,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventRequestCompleted,
		constants.KeyHotelName, hotel,
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
			constants.KeyHotelName, hotel,
		)

		http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		return
	}

	logger.Info(
		constants.EventRequestCompleted,
		constants.KeyHandler, "GetByHotel",
		constants.KeyHotelName, hotel,
		constants.KeyStatus, http.StatusOK,
	)
}
