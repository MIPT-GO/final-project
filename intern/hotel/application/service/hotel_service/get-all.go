package hotel_service

import (
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
