package hotel_service

import (
	"errors"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
)

func (s *HotelServiceImpl) GetEmail(name string) (string, error) {
	s.Log.Info(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventHotelGetEmail,
		logs.KeyHotelName, name)

	email, err := s.HotelRep.GetByEmail(name)

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
		logs.KeyHotelEmail, email)

	return email, nil
}
