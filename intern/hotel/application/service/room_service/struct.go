package room_service

import (
	"final-project/intern/hotel/domain/interfaces"
	"log/slog"
)

type RoomServiceImpl struct {
	Log     *slog.Logger
	RoomRep interfaces.RoomRepository
}

func NewRoomService(
	log *slog.Logger,
	roomRep interfaces.RoomRepository,
) interfaces.RoomService {
	return &RoomServiceImpl{Log: log, RoomRep: roomRep}
}
