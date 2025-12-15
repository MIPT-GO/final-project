package service

import (
	"fmt"

	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/application/operations"
	"final-project/intern/booking/domain/models/reservation"
	"final-project/pkg/booking/constants"
)

type ReservationService struct {
	repo     interfaces.Repository
	webhook  string
	Logger   interfaces.Logger
	Producer interfaces.Producer
}

func NewReservationService(repo interfaces.Repository, webhook string, logger interfaces.Logger, producer interfaces.Producer) ReservationService {
	return ReservationService{repo: repo, webhook: webhook, Logger: logger, Producer: producer}
}

func (s *ReservationService) GetById(id uint64) (reservation.Reserve, error) {
	if s.repo == nil {
		err := constants.ErrRepoNotSpecified
		s.Logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch)
		return reservation.Reserve{}, err
	}
	res, err := s.repo.GetById(id)
	if err != nil {
		s.Logger.Error(constants.EventFailedToFetch, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventFailedToFetch)
		return reservation.Reserve{}, err
	}
	return res, nil
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

func (s *ReservationService) BookRoomInHotel(reserve reservation.Reserve) (uint64, string, error) {
	s.Logger.Info(constants.EventBookRoomInHotel, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel, constants.KeyHotelName, reserve.Hotel, constants.KeyRoomNumber, reserve.Number, constants.KeyUserEmail, reserve.Email)
	if s.repo == nil {
		err := constants.ErrRepoNotSpecified
		s.Logger.Error(constants.EventBookRoomInHotel, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel)
		return 0, "", err
	}

	id, err := operations.BookRoomInHotel(reserve, s.repo)
	if err != nil {
		s.Logger.Error(constants.EventBookRoomInHotel, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel, constants.KeyHotelName, reserve.Hotel, constants.KeyRoomNumber, reserve.Number)
		return 0, "", err
	}

	amount, err := s.repo.GetRoomPrice(reserve.Hotel, reserve.Number)
	if err != nil {
		s.Logger.Error(constants.EventHotelRequest, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest)
		return 0, "", err
	}

	// webhook will be called with reservation id only
	webhookWithID := s.webhook + "/" + fmt.Sprintf("%d", id)

	link, err := s.repo.InitiatePayment(amount, webhookWithID, "Reservation payment", map[string]string{"id": fmt.Sprintf("%d", id)})
	if err != nil {
		s.Logger.Error(constants.EventPaymentInitiated, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventPaymentInitiated)
		return 0, "", err
	}

	s.Logger.Info(constants.EventBookRoomInHotel, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventBookRoomInHotel, constants.KeyHotelName, reserve.Hotel, constants.KeyRoomNumber, reserve.Number, constants.KeyStatus, "booked")

	return id, link, nil
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
