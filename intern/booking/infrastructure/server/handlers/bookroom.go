package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"final-project/intern/booking/constants"
	"final-project/intern/booking/domain/models/reservation"
	"final-project/intern/booking/infrastructure/server/dto"
)

func (handl *handler) BookRoomInHotel(w http.ResponseWriter, r *http.Request) {
	logger := handl.service.Logger

	logger.Info(
		constants.EventIncomingRequest,
		constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventIncomingRequest,
		constants.KeyHandler, "BookRoomInHotel",
		constants.KeyMethod, r.Method,
		constants.KeyPath, r.URL.Path,
	)

	if r.Method != http.MethodPost {
		http.Error(w, constants.MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	var req dto.Reservation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("failed to decode request body", constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyError, err)
		http.Error(w, constants.MsgUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}

	start, _ := time.Parse(time.RFC3339, req.Start)
	end, _ := time.Parse(time.RFC3339, req.End)
	reservation := reservation.Reserve{
		Email:  req.Email,
		Start:  start,
		End:    end,
		Hotel:  req.Hotel,
		Number: req.Number,
	}
	err := handl.service.CheckAccuracy(reservation)
	if err != nil {
		if errors.Is(err, constants.ErrHotelOrRoomNotFound) {
			logger.Warn(constants.EventFailedToFetch, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyHotelName, reservation.Hotel, constants.KeyRoomNumber, reservation.Number)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if errors.Is(err, constants.ErrRoomAlreadyBooked) {
			logger.Warn(constants.EventFailedToFetch, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyHotelName, reservation.Hotel, constants.KeyRoomNumber, reservation.Number)
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}

		logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyHandler, "BookRoomInHotel")
		http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		return
	}

	logger.Info(constants.EventBookRoomInHotel, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel, constants.KeyHandler, "BookRoomInHotel", constants.KeyHotelName, reservation.Hotel, constants.KeyRoomNumber, reservation.Number, constants.KeyUserEmail, reservation.Email)
	id, paymentLink, err := handl.service.BookRoomInHotel(reservation)
	if err != nil {
		logger.Error(constants.EventBookRoomInHotel, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel, constants.KeyHandler, "BookRoomInHotel")
		http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"id": fmt.Sprintf("%d", id), "payment_link": paymentLink}
	logger.Info(constants.MsgPaymentLinkCreated, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventPaymentInitiated, constants.KeyPaymentLink, paymentLink)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error(constants.EventFailedToEncode, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToEncode)
	}
}
