package room_service

import (
	"final-project/intern/hotel/domain/interfaces"
	"log/slog"
)

type RoomServiceImpl struct {
	Log      *slog.Logger
	HotelRep interfaces.HotelRepository
	RoomRep  interfaces.RoomRepository
}

func NewRoomService(
	log *slog.Logger,
	hotelRep interfaces.HotelRepository,
	roomRep interfaces.RoomRepository,
) interfaces.RoomService {
	return &RoomServiceImpl{Log: log, HotelRep: hotelRep, RoomRep: roomRep}
}
