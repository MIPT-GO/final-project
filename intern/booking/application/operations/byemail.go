package operations

import (
	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/domain/models/reservation"
)

func GetByEmail(email string, repo interfaces.Repository) ([]reservation.Reserve, error) {
	reservations, err := repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}
