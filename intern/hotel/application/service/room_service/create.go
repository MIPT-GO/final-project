package room_service

import (
	"errors"
	"final-project/intern/hotel/domain/entity"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
)

func (s *RoomServiceImpl) Create(hotelName string, number int, cost float32, userEmail string) error {

	s.Log.Info(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventRoomCreate,
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
			logs.KeyError, err)
		return custom_errors.ErrDatabaseFailure
	}

	if ownerEmail != userEmail {
		s.Log.Warn(logs.MsgPermissionDenied,
			logs.KeyEvent, logs.EventPermissionDenied,
			logs.KeyUserEmail, userEmail)
		return custom_errors.ErrPermissionDenied
	}

	newRoom := entity.Room{
		HotelName: hotelName,
		Number:    number,
		Cost:      cost,
	}
	err = s.RoomRep.AddNewRoom(newRoom)
	if err != nil {
		if errors.Is(err, custom_errors.ErrEntityAlreadyExists) {
			s.Log.Error(logs.MsgEntityAlreadyExists,
				logs.KeyHotelName, hotelName,
				logs.KeyRoomNumber, number)
			return custom_errors.ErrEntityAlreadyExists
		}
		s.Log.Error(logs.MsgDatabaseWriteFailed,
			logs.KeyError, err)
		return custom_errors.ErrDatabaseFailure
	}

	s.Log.Info(logs.MsgOperationSuccess, logs.KeyHotelName, hotelName)
	return nil
}
