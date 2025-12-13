package hotel_service

import (
	"errors"
	"final-project/intern/hotel/domain/entity"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
)

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
