package hotel_service

import (
	"final-project/intern/hotel/domain/interfaces"
	"log/slog"
)

type HotelServiceImpl struct {
	Log      *slog.Logger
	HotelRep interfaces.HotelRepository
}

func NewHotelService(log *slog.Logger, hotelRep interfaces.HotelRepository) interfaces.HotelService {
	return &HotelServiceImpl{Log: log, HotelRep: hotelRep}
}
