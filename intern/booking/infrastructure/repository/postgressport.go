package repository

import "final-project/intern/booking/domain/models/reservation"

type PostgresPort interface {
	FindByEmail(email string) ([]reservation.Reserve, error)
	FindByHotel(hotel string) ([]reservation.Reserve, error)
	GetById(id uint64) (reservation.Reserve, error)
	AddNewReservation(reserv reservation.Reserve) (uint64, error)
	HasOverlap(reserv reservation.Reserve) (bool, error)
	DeleteReservation(reserv reservation.Reserve) error
}
