package handlers

import (
	"encoding/json"
	"final-project/intern/hotel/domain/interfaces"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if body != nil {
		json.NewEncoder(w).Encode(body)
	}
}

func getPathParam(r *http.Request, name string) string {
	vars := mux.Vars(r)
	return vars[name]
}

type HotelHandlerImpl struct {
	Log          *slog.Logger
	HotelService interfaces.HotelService
	RoomService  interfaces.RoomService
}

func NewHotelHandler(log *slog.Logger, hs interfaces.HotelService, rs interfaces.RoomService) interfaces.HotelHandler {
	return &HotelHandlerImpl{
		Log:          log,
		HotelService: hs,
		RoomService:  rs,
	}
}
