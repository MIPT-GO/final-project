package operations

import (
	"final-project/intern/booking/application/interfaces"
)

func GetAvailableInHotel(hotel string, repo interfaces.Repository) ([]uint64, error) {
	rooms, err := repo.GetAvailableInHotel(hotel)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
