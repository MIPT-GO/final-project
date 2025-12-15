package mock

import (
	"final-project/intern/hotel/domain/entity"

	"github.com/stretchr/testify/mock"
)

type HotelRepositoryMock struct {
	mock.Mock
}

func (m *HotelRepositoryMock) GetAll() ([]entity.Hotel, error) {
	args := m.Called()

	var result []entity.Hotel
	if args.Get(0) != nil {
		result = args.Get(0).([]entity.Hotel)
	}

	return result, args.Error(1)
}

func (m *HotelRepositoryMock) GetByName(name string) (entity.Hotel, error) {
	args := m.Called(name)

	if args.Get(0) == nil {
		return entity.Hotel{}, args.Error(1)
	}

	return args.Get(0).(entity.Hotel), args.Error(1)
}

func (m *HotelRepositoryMock) GetByEmail(email string) (string, error) {
	args := m.Called(email)

	if args.Get(0) == nil {
		return "", args.Error(1)
	}

	return args.String(0), args.Error(1)
}

func (m *HotelRepositoryMock) AddNewHotel(hotel entity.Hotel) error {
	args := m.Called(hotel)
	return args.Error(0)
}

func (m *HotelRepositoryMock) ExistsByName(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}
