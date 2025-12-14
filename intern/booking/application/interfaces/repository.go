package interfaces

import "final-project/intern/booking/domain/models/reservation"

type Repository interface {
	FindByEmail(email string) ([]reservation.Reserve, error)
	FindByHotel(hotel string) ([]reservation.Reserve, error)
	CheckAccuracy(reserv reservation.Reserve) error
	AddNewReservation(reserv reservation.Reserve) error
	DeleteReservation(reserv reservation.Reserve) error
	GetRoomPrice(hotel string, roomNumber uint64) (string, error)
	GetAvailableInHotel(hotel string) ([]uint64, error)
	InitiatePayment(amount string, webhook string, message string, extra map[string]string) (string, error)
}
