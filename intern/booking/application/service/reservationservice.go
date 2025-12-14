package service

import (
	"fmt"
	"time"

	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/application/operations"
	"final-project/intern/booking/domain/models/reservation"
	"final-project/pkg/booking/constants"
)

type ReservationService struct {
	repo    interfaces.Repository
	webhook string
	Logger  interfaces.Logger
}

func NewReservationService(repo interfaces.Repository, webhook string, logger interfaces.Logger) ReservationService {
	return ReservationService{repo: repo, webhook: webhook, Logger: logger}
}

func (s *ReservationService) GetByEmail(email string) ([]reservation.Reserve, error) {
	s.Logger.Debug(constants.EventGetByEmail, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventGetByEmail, constants.KeyUserEmail, email)
	if s.repo == nil {
		err := constants.ErrRepoNotSpecified
		s.Logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch)
		return nil, err
	}

	reservs, err := operations.GetByEmail(email, s.repo)
	if err != nil {
		s.Logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyUserEmail, email)
		return nil, err
	}
	s.Logger.Info(constants.EventRequestCompleted, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventRequestCompleted, constants.KeyUserEmail, email, constants.KeyCount, len(reservs))
	return reservs, nil
}

func (s *ReservationService) GetByHotel(hotel string) ([]reservation.Reserve, error) {
	s.Logger.Debug(constants.EventGetByHotel, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventGetByHotel, constants.KeyHotelName, hotel)
	if s.repo == nil {
		err := constants.ErrRepoNotSpecified
		s.Logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch)
		return nil, err
	}

	reservs, err := operations.GetByHotel(hotel, s.repo)
	if err != nil {
		s.Logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyHotelName, hotel)
		return nil, err
	}
	s.Logger.Info(constants.EventRequestCompleted, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventRequestCompleted, constants.KeyHotelName, hotel, constants.KeyCount, len(reservs))
	return reservs, nil
}

func (s *ReservationService) BookRoomInHotel(reserve reservation.Reserve) (string, error) {
	s.Logger.Info(constants.EventBookRoomInHotel, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel, constants.KeyHotelName, reserve.Hotel, constants.KeyRoomNumber, reserve.Number, constants.KeyUserEmail, reserve.Email)
	if s.repo == nil {
		err := constants.ErrRepoNotSpecified
		s.Logger.Error(constants.EventBookRoomInHotel, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel)
		return "", err
	}

	if err := operations.BookRoomInHotel(reserve, s.repo); err != nil {
		s.Logger.Error(constants.EventBookRoomInHotel, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel, constants.KeyHotelName, reserve.Hotel, constants.KeyRoomNumber, reserve.Number)
		return "", err
	}
	extra := map[string]string{
		"email":  reserve.Email,
		"hotel":  reserve.Hotel,
		"number": fmt.Sprintf("%d", reserve.Number),
		"start":  reserve.Start.Format(time.RFC3339),
		"end":    reserve.End.Format(time.RFC3339),
	}

	amount, err := s.repo.GetRoomPrice(reserve.Hotel, reserve.Number)
	if err != nil {
		s.Logger.Error(constants.EventHotelRequest, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest)
		return "", err
	}

	link, err := s.repo.InitiatePayment(amount, s.webhook, "Reservation payment", extra)
	if err != nil {
		s.Logger.Error(constants.EventPaymentInitiated, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventPaymentInitiated)
		return "", err
	}

	s.Logger.Info(constants.EventBookRoomInHotel, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel, constants.KeyHotelName, reserve.Hotel, constants.KeyRoomNumber, reserve.Number, constants.KeyStatus, "booked")
	return link, nil
}

func (s *ReservationService) CheckAccuracy(reserve reservation.Reserve) error {
	s.Logger.Debug(constants.EventCheckAccuracy, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventCheckAccuracy, constants.KeyHotelName, reserve.Hotel, constants.KeyRoomNumber, reserve.Number)
	if s.repo == nil {
		err := constants.ErrRepoNotSpecified
		s.Logger.Error(constants.EventCheckAccuracy, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventCheckAccuracy)
		return err
	}

	if err := operations.CheckAccuracy(reserve, s.repo); err != nil {
		s.Logger.Warn(constants.EventCheckAccuracy, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventCheckAccuracy, constants.KeyError, err, constants.KeyHotelName, reserve.Hotel, constants.KeyRoomNumber, reserve.Number)
		return err
	}
	s.Logger.Info(constants.EventCheckAccuracy, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventCheckAccuracy, constants.KeyHotelName, reserve.Hotel, constants.KeyRoomNumber, reserve.Number, constants.KeyStatus, "ok")
	return nil
}

func (s *ReservationService) GetAvailableInHotel(hotel string) ([]uint64, error) {
	s.Logger.Debug(constants.EventGetAvailableNumbers, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventGetAvailableNumbers, constants.KeyHotelName, hotel)
	if s.repo == nil {
		err := constants.ErrRepoNotSpecified
		s.Logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch)
		return nil, err
	}

	rooms, err := operations.GetAvailableInHotel(hotel, s.repo)
	if err != nil {
		s.Logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch, constants.KeyHotelName, hotel)
		return nil, err
	}
	s.Logger.Info(constants.EventRequestCompleted, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventRequestCompleted, constants.KeyHotelName, hotel, constants.KeyCount, len(rooms))
	return rooms, nil
}

func (s *ReservationService) DeleteReservation(reserv reservation.Reserve) error {
	s.Logger.Info(constants.EventDeleteReservation, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDeleteReservation, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number, constants.KeyUserEmail, reserv.Email)
	if s.repo == nil {
		err := constants.ErrRepoNotSpecified
		s.Logger.Error(constants.EventDeleteReservation, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDeleteReservation)
		return err
	}

	if err := operations.DeleteReservation(reserv, s.repo); err != nil {
		s.Logger.Error(constants.EventDeleteReservation, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDeleteReservation, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number)
		return err
	}
	s.Logger.Info(constants.EventDeleteReservation, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDeleteReservation, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number, constants.KeyStatus, "deleted")
	return nil
}
