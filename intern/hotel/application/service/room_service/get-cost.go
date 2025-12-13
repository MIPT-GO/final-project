package room_service

import (
	"errors"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
)

func (s *RoomServiceImpl) GetCost(hotelName string, number int) (float32, error) {
	s.Log.Info(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventRoomCostQuery,
		logs.KeyHotelName, hotelName,
		logs.KeyRoomNumber, number)

	cost, err := s.RoomRep.GetCost(hotelName, number)

	if errors.Is(err, custom_errors.ErrEntityNotFound) {
		s.Log.Warn(logs.MsgEntityNotFound,
			logs.KeyHotelName, hotelName,
			logs.KeyRoomNumber, number)
		return 0, custom_errors.ErrEntityNotFound
	}
	if err != nil {
		s.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyHotelName, hotelName,
			logs.KeyError, err)
		return 0, custom_errors.ErrDatabaseFailure
	}

	s.Log.Info(logs.MsgOperationSuccess,
		logs.KeyHotelName, hotelName,
		logs.KeyRoomNumber, number,
		logs.KeyCost, cost)
	return cost, nil
}
