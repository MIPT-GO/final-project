package operations

import (
	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/domain/models/reservation"
)

func BookRoomInHotel(reserv reservation.Reserve, repo interfaces.Repository) error {
	err := repo.AddNewReservation(reserv)
	return err
}
