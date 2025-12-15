package room_service

import (
	"errors"
	"final-project/intern/hotel/domain/entity"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
	"log/slog"
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
