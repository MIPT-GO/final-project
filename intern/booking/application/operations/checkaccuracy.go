package operations

import (
	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/domain/models/reservation"
)

func CheckAccuracy(reserv reservation.Reserve, repo interfaces.Repository) error {
	return repo.CheckAccuracy(reserv)
}
