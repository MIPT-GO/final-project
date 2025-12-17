package repository

import (
	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/constants"
	"final-project/intern/booking/domain/models/reservation"
	"time"
)

type ComplexRepo struct {
	postgres PostgresPort
	hotel    HotelPort
	payment  PaymentPort
	logger   interfaces.Logger
}

func NewComplexRepo(pg PostgresPort, ht HotelPort, pay PaymentPort, logger interfaces.Logger) *ComplexRepo {
	if logger != nil {
		logger.Info(constants.EventRepoInit, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventRepoInit)
	}
	return &ComplexRepo{
		postgres: pg,
		hotel:    ht,
		payment:  pay,
		logger:   logger,
	}
}

func (cr *ComplexRepo) FindByEmail(email string) ([]reservation.Reserve, error) {
	if cr.logger != nil {
		cr.logger.Debug(constants.EventDBQuery, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBQuery, constants.KeyUserEmail, email)
	}
	return cr.postgres.FindByEmail(email)
}

func (cr *ComplexRepo) FindByHotel(hotel string) ([]reservation.Reserve, error) {
	if cr.logger != nil {
		cr.logger.Debug(constants.EventDBQuery, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBQuery, constants.KeyHotelName, hotel)
	}
	return cr.postgres.FindByHotel(hotel)
}

func (cr *ComplexRepo) DeleteReservation(reserv reservation.Reserve) error {
	if cr.logger != nil {
		cr.logger.Info(constants.EventDBDelete, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBDelete, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number)
	}
	return cr.postgres.DeleteReservation(reserv)
}

func (cr *ComplexRepo) AddNewReservation(reserv reservation.Reserve) (uint64, error) {
	if cr.logger != nil {
		cr.logger.Info(constants.EventDBInsert, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBInsert, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number)
	}
	return cr.postgres.AddNewReservation(reserv)
}

func (cr *ComplexRepo) GetById(id uint64) (reservation.Reserve, error) {
	if cr.logger != nil {
		cr.logger.Debug(constants.EventDBQuery, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBQuery, constants.KeyQuery, constants.QueryFindById, constants.KeyUserID, id)
	}
	return cr.postgres.GetById(id)
}

func (cr *ComplexRepo) InitiatePayment(amount string, webhook string, message string, extra map[string]string) (string, error) {
	if cr.logger != nil {
		cr.logger.Info(constants.EventPaymentInitiated, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventPaymentInitiated, constants.KeyAmount, amount)
	}
	return cr.payment.InitiatePayment(amount, webhook, message, extra)
}

func (cr *ComplexRepo) GetRoomPrice(hotel string, roomNumber uint64) (string, error) {
	if cr.logger != nil {
		cr.logger.Debug(constants.EventHotelRequest, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyHotelName, hotel, constants.KeyRoomNumber, roomNumber)
	}
	return cr.hotel.GetRoomPrice(hotel, roomNumber)
}

func (cr *ComplexRepo) GetHotelOwnerEmail(hotel string) (string, error) {
	if cr.logger != nil {
		cr.logger.Debug(constants.EventHotelRequest, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyHotelName, hotel)
	}
	return cr.hotel.GetOwnerEmail(hotel)
}

func (cr *ComplexRepo) CheckAccuracy(reserv reservation.Reserve) error {
	if cr.logger != nil {
		cr.logger.Debug(constants.EventHasOverlap, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHasOverlap, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number)
	}
	overlap, err := cr.postgres.HasOverlap(reserv)
	if err != nil {
		if cr.logger != nil {
			cr.logger.Error(constants.EventHasOverlap, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHasOverlap)
		}
		return err
	}
	if overlap {
		if cr.logger != nil {
			cr.logger.Warn(constants.EventHasOverlap, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHasOverlap, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number)
		}
		return constants.ErrRoomAlreadyBooked
	}

	if cr.logger != nil {
		cr.logger.Debug(constants.EventHotelRequest, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number)
	}
	correct, err := cr.hotel.CheckAccuracy(reserv.Hotel, reserv.Number)
	if err != nil {
		if cr.logger != nil {
			cr.logger.Error(constants.EventHotelRequest, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest)
		}
		return err
	}
	if !correct {
		if cr.logger != nil {
			cr.logger.Warn(constants.EventHotelRequest, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyHotelName, reserv.Hotel)
		}
		return constants.ErrHotelOrRoomNotFound
	}
	return nil
}

func (cr *ComplexRepo) GetAvailableInHotel(hotel string) ([]uint64, error) {
	rooms, err := cr.hotel.GetAllRoomsInHotel(hotel)
	if err != nil {
		return nil, err
	}
	availableRooms := make([]uint64, 0)
	for _, room := range rooms {
		reserv := reservation.Reserve{
			Hotel:  hotel,
			Number: room,
			Start:  time.Now(),
			End:    time.Now(),
		}
		overlap, err := cr.postgres.HasOverlap(reserv)
		if err != nil {
			return nil, err
		}
		if !overlap {
			availableRooms = append(availableRooms, room)
		}
	}
	return availableRooms, nil
}
