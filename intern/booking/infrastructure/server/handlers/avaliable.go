package handlers

import (
	"encoding/json"
	"errors"
	"final-project/intern/booking/infrastructure/server/dto"
	"final-project/pkg/booking/constants"
	"net/http"
)

func (handl *handler) Available(w http.ResponseWriter, r *http.Request) {
	log := handl.service.Logger

	log.Info(
		constants.EventIncomingRequest,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventIncomingRequest,
		constants.KeyHandler, "Available",
		constants.KeyMethod, r.Method,
		constants.KeyPath, r.URL.Path,
	)

	if r.Method != http.MethodGet {
		log.Warn(
			constants.EventMethodNotAllowed,
			constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventMethodNotAllowed,
			constants.KeyHandler, "Available",
			constants.KeyMethod, r.Method,
		)

		http.Error(w, constants.MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	hotel := r.PathValue("hotel")

	log.Debug(
		constants.EventGetAvailableNumbers,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventGetAvailableNumbers,
		constants.KeyHotelName, hotel,
	)

	numbers, err := handl.service.GetAvailableInHotel(hotel)
	if err != nil {
		if errors.Is(err, constants.ErrHotelNotFound) {
			log.Warn(
				constants.EventFailedToFetch,
				constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch,
				constants.KeyHotelName, hotel,
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Error(
			constants.EventFailedToFetch,
			err,
			constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch,
			constants.KeyHotelName, hotel,
		)

		http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		return
	}

	log.Info(
		constants.EventRequestCompleted,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventRequestCompleted,
		constants.KeyHotelName, hotel,
		constants.KeyCount, len(numbers),
	)

	out := dto.Numbers{
		Numbers: numbers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Error(
			constants.EventFailedToEncode,
			err,
			constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToEncode,
			constants.KeyHotelName, hotel,
		)

		http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		return
	}

	log.Info(
		constants.EventRequestCompleted,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventRequestCompleted,
		constants.KeyHandler, "Available",
		constants.KeyHotelName, hotel,
		constants.KeyStatus, http.StatusOK,
	)
}
