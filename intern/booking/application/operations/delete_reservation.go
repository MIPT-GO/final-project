package operations

import (
	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/domain/models/reservation"
)

func DeleteReservation(reserv reservation.Reserve, repo interfaces.Repository) error {
	return repo.DeleteReservation(reserv)
}
