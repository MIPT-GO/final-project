package room_service

import (
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
)

func (s *RoomServiceImpl) GetAll(hotelName string) ([]int, error) {
	s.Log.Info(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventRoomCostQuery,
		logs.KeyHotelName, hotelName)

	rooms, err := s.RoomRep.GetAll(hotelName)

	if err != nil {
		s.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyHotelName, hotelName,
			logs.KeyError, err)
		return nil, custom_errors.ErrDatabaseFailure
	}

	roomNumbers := make([]int, 0, len(rooms))
	for _, room := range rooms {
		roomNumbers = append(roomNumbers, room.Number)
	}

	s.Log.Info(logs.MsgOperationSuccess,
		logs.KeyHotelName, hotelName,
		"count", len(roomNumbers))

	return roomNumbers, nil
}
