package hotel_service

import (
	"errors"
	"final-project/intern/hotel/domain/entity"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
)

func (s *HotelServiceImpl) GetAll() ([]string, error) {
	s.Log.Info(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventHotelGetAll)

	hotelEntities, err := s.HotelRep.GetAll()

	if err != nil {
		s.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventHotelGetAll,
			logs.KeyError, err)
		return nil, custom_errors.ErrDatabaseFailure
	}

	hotelNames := make([]string, 0, len(hotelEntities))
	for _, hotel := range hotelEntities {
		hotelNames = append(hotelNames, hotel.Name)
	}

	s.Log.Info(logs.MsgOperationSuccess,
		logs.KeyEvent, logs.EventHotelGetAll,
		"count", len(hotelNames))

	return hotelNames, nil
}

func (s *HotelServiceImpl) GetEmail(name string) (string, error) {
	s.Log.Info(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventHotelGetEmail,
		logs.KeyHotelName, name)

	hotel, err := s.HotelRep.GetByName(name)

	if errors.Is(err, custom_errors.ErrEntityNotFound) {
		s.Log.Warn(logs.MsgEntityNotFound,
			logs.KeyHotelName, name)
		return "", custom_errors.ErrEntityNotFound
	}
	if err != nil {
		s.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventHotelGetEmail,
			logs.KeyError, err)
		return "", custom_errors.ErrDatabaseFailure
	}

	s.Log.Info(logs.MsgOperationSuccess,
		logs.KeyHotelName, name,
		logs.KeyHotelEmail, hotel.Email)

	return hotel.Email, nil
}

func (s *HotelServiceImpl) Create(name string, email string) error {
	s.Log.Info(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventHotelCreate,
		logs.KeyHotelName, name,
		logs.KeyHotelEmail, email)

	newHotel := entity.Hotel{
		Name:  name,
		Email: email,
	}

	err := s.HotelRep.AddNewHotel(newHotel)

	if err != nil {
		if errors.Is(err, custom_errors.ErrEntityAlreadyExists) {
			s.Log.Error(logs.MsgEntityAlreadyExists,
				logs.KeyHotelName, name)
			return custom_errors.ErrEntityAlreadyExists
		}
		s.Log.Error(logs.MsgDatabaseWriteFailed,
			logs.KeyEvent, logs.EventHotelCreate,
			logs.KeyError, err)
		return custom_errors.ErrDatabaseFailure
	}

	s.Log.Info(logs.MsgOperationSuccess, logs.KeyHotelName, name)
	return nil
}
