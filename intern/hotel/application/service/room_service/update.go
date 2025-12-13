package room_service

import (
	"errors"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
	"log/slog"
)

func (s *RoomServiceImpl) UpdateCost(hotelName string, number int, newCost float32, userEmail string) error {
	s.Log.Info(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventRoomCostUpdate,
		logs.KeyHotelName, hotelName,
		logs.KeyRoomNumber, number,
		logs.KeyUserEmail, userEmail)

	ownerEmail, err := s.HotelRep.GetByEmail(hotelName)

	if errors.Is(err, custom_errors.ErrEntityNotFound) {
		s.Log.Warn(logs.MsgEntityNotFound,
			logs.KeyHotelName, hotelName)
		return custom_errors.ErrEntityNotFound
	}
	if err != nil {
		s.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyHotelName, hotelName,
			logs.KeyError, err)
		return custom_errors.ErrDatabaseFailure
	}

	if ownerEmail != userEmail {
		s.Log.Warn(logs.MsgPermissionDenied,
			logs.KeyEvent, logs.EventPermissionDenied,
			logs.KeyUserEmail, userEmail)
		return custom_errors.ErrPermissionDenied
	}

	err = s.RoomRep.UpdateCost(hotelName, number, newCost)

	if errors.Is(err, custom_errors.ErrEntityNotFound) {
		s.Log.Warn(logs.MsgUpdateFailure,
			logs.MsgEntityNotFound,
			logs.KeyHotelName, hotelName,
			slog.Int(logs.KeyRoomNumber, number))
		return custom_errors.ErrEntityNotFound
	}
	if err != nil {
		s.Log.Error(logs.MsgDatabaseWriteFailed,
			logs.KeyHotelName, hotelName,
			logs.KeyError, err)
		return custom_errors.ErrDatabaseFailure
	}

	s.Log.Info(logs.MsgOperationSuccess,
		logs.KeyHotelName, hotelName,
		logs.KeyRoomNumber, number,
		logs.KeyNewCost, newCost)
	return nil
}
