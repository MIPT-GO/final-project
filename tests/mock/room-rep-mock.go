package mock

import (
	"final-project/intern/hotel/domain/entity"

	"github.com/stretchr/testify/mock"
)

type RoomRepositoryMock struct {
	mock.Mock
}

func (m *RoomRepositoryMock) GetAll(hotelName string) ([]entity.Room, error) {
	args := m.Called(hotelName)

	var result []entity.Room
	if args.Get(0) != nil {
		result = args.Get(0).([]entity.Room)
	}

	return result, args.Error(1)
}

func (m *RoomRepositoryMock) GetCost(hotelName string, number int) (float32, error) {
	args := m.Called(hotelName, number)
	return float32(args.Get(0).(float32)), args.Error(1)
}

func (m *RoomRepositoryMock) AddNewRoom(room entity.Room) error {
	args := m.Called(room)
	return args.Error(0)
}

func (m *RoomRepositoryMock) UpdateCost(hotelName string, number int, newCost float32) error {
	args := m.Called(hotelName, number, newCost)
	return args.Error(0)
}
