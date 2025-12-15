package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"final-project/intern/hotel/infrastructure/server/dto"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
)

// GET /v1/hotel/all/
func (h *HotelHandlerImpl) GetAllHotels(w http.ResponseWriter, r *http.Request) {
	hotels, err := h.HotelService.GetAll()

	if err != nil {
		h.Log.Error(logs.MsgDatabaseQueryFailed, logs.KeyEvent, logs.EventHotelGetAll, logs.KeyError, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := dto.AllHotelsResponse{Hotels: hotels}
	RespondJSON(w, http.StatusOK, response)
}

// POST /v1/hotel/
func (h *HotelHandlerImpl) CreateHotel(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateHotelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Warn("Failed to decode request body", logs.KeyError, err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	err := h.HotelService.Create(req.Name, req.Email)

	if errors.Is(err, custom_errors.ErrEntityAlreadyExists) {
		h.Log.Warn(logs.MsgEntityAlreadyExists, logs.KeyHotelName, req.Name)
		http.Error(w, "Hotel already exists", http.StatusConflict)
		return
	}
	if errors.Is(err, custom_errors.ErrDatabaseFailure) {
		h.Log.Error(logs.MsgDatabaseWriteFailed, logs.KeyEvent, logs.EventHotelCreate, logs.KeyError, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	RespondJSON(w, http.StatusCreated, nil)
}

// GET /v1/hotel/{hotel: str}/
func (h *HotelHandlerImpl) GetOneHotel(w http.ResponseWriter, r *http.Request) {
	hotelName := getPathParam(r, "hotel")

	email, err := h.HotelService.GetEmail(hotelName)

	if errors.Is(err, custom_errors.ErrEntityNotFound) {
		h.Log.Warn(logs.MsgEntityNotFound, logs.KeyHotelName, hotelName)
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, custom_errors.ErrDatabaseFailure) {
		h.Log.Error(logs.MsgDatabaseQueryFailed, logs.KeyHotelName, hotelName, logs.KeyError, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := dto.GetHotelResponse{Name: hotelName, Email: email}
	RespondJSON(w, http.StatusOK, response)
}

// GET /v1/hotel/{hotel: str}/rooms/
func (h *HotelHandlerImpl) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	hotelName := getPathParam(r, "hotel")

	rooms, err := h.RoomService.GetAll(hotelName)

	if err != nil {
		h.Log.Error(logs.MsgDatabaseQueryFailed, logs.KeyEvent, logs.EventRoomCostQuery, logs.KeyHotelName, hotelName, logs.KeyError, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := dto.AllRoomsResponse{Numbers: rooms}
	RespondJSON(w, http.StatusOK, response)
}

// PUT /v1/hotel/{hotel: str}/rooms/
func (h *HotelHandlerImpl) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	hotelName := getPathParam(r, "hotel")

	var req dto.UpdateRoomRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Warn("Failed to decode UpdateRoom request", logs.KeyError, err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	err := h.RoomService.UpdateCost(hotelName, req.Number, req.NewCost)

	if errors.Is(err, custom_errors.ErrEntityNotFound) {
		h.Log.Warn(logs.MsgEntityNotFound, logs.KeyHotelName, hotelName, slog.Int(logs.KeyRoomNumber, req.Number))
		http.Error(w, "Room or Hotel not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, custom_errors.ErrDatabaseFailure) {
		h.Log.Error(logs.MsgDatabaseWriteFailed, logs.KeyHotelName, hotelName, logs.KeyError, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	RespondJSON(w, http.StatusOK, nil)
}

// GET /v1/hotel/{hotel}/rooms/{number}
func (h *HotelHandlerImpl) GetRoom(w http.ResponseWriter, r *http.Request) {
	hotelName := getPathParam(r, "hotel")
	roomStr := getPathParam(r, "number")

	roomNumber, err := strconv.Atoi(roomStr)
	if err != nil {
		h.Log.Warn("Invalid room number", logs.KeyError, err)
		http.Error(w, "Invalid room number", http.StatusBadRequest)
		return
	}

	cost, err := h.RoomService.GetCost(hotelName, roomNumber)

	if errors.Is(err, custom_errors.ErrEntityNotFound) {
		h.Log.Warn(
			logs.MsgEntityNotFound,
			logs.KeyHotelName, hotelName,
			slog.Int(logs.KeyRoomNumber, roomNumber),
		)
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if errors.Is(err, custom_errors.ErrDatabaseFailure) {
		h.Log.Error(
			logs.MsgDatabaseQueryFailed,
			logs.KeyHotelName, hotelName,
			slog.Int(logs.KeyRoomNumber, roomNumber),
			logs.KeyError, err,
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := dto.GetRoomResponse{
		Cost: cost,
	}

	RespondJSON(w, http.StatusOK, response)
}

// POST /v1/hotel/{hotel}/rooms/
func (h *HotelHandlerImpl) CreateRoom(w http.ResponseWriter, r *http.Request) {
	hotelName := getPathParam(r, "hotel")

	var req dto.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Warn("Failed to decode CreateRoom request", logs.KeyError, err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	err := h.RoomService.Create(hotelName, req.Number, req.Cost)

	if errors.Is(err, custom_errors.ErrEntityAlreadyExists) {
		h.Log.Warn(logs.MsgEntityAlreadyExists, logs.KeyHotelName, hotelName, slog.Int(logs.KeyRoomNumber, req.Number))
		http.Error(w, "Room already exists", http.StatusConflict)
		return
	}
	if errors.Is(err, custom_errors.ErrEntityNotFound) {
		h.Log.Warn(logs.MsgEntityNotFound, logs.KeyHotelName, hotelName)
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, custom_errors.ErrDatabaseFailure) {
		h.Log.Error(logs.MsgDatabaseWriteFailed, logs.KeyHotelName, hotelName, logs.KeyError, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	RespondJSON(w, http.StatusCreated, nil)
}
