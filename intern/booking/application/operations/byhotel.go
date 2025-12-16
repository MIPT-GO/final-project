package operations

import (
	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/domain/models/reservation"
)

func GetByHotel(hotel string, repo interfaces.Repository) ([]reservation.Reserve, error) {
	reservations, err := repo.FindByHotel(hotel)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}
