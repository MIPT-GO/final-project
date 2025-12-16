package mock

import (
	"final-project/intern/hotel/domain/interfaces"
)

type MockHotelService struct {
	GetAllFunc   func() ([]string, error)
	GetEmailFunc func(name string) (string, error)
	CreateFunc   func(name, email string) error
}

func (m *MockHotelService) GetAll() ([]string, error)            { return m.GetAllFunc() }
func (m *MockHotelService) GetEmail(name string) (string, error) { return m.GetEmailFunc(name) }
func (m *MockHotelService) Create(name, email string) error      { return m.CreateFunc(name, email) }

type MockRoomService struct {
	GetAllFunc     func(name string) ([]int, error)
	CreateFunc     func(hotelName string, number int, cost float32) error
	GetCostFunc    func(name string, number int) (float32, error)
	UpdateCostFunc func(hotelName string, number int, newCost float32) error
}

func (m *MockRoomService) GetAll(name string) ([]int, error) { return m.GetAllFunc(name) }
func (m *MockRoomService) Create(hotelName string, number int, cost float32) error {
	return m.CreateFunc(hotelName, number, cost)
}
func (m *MockRoomService) GetCost(name string, number int) (float32, error) {
	return m.GetCostFunc(name, number)
}
func (m *MockRoomService) UpdateCost(hotelName string, number int, newCost float32) error {
	return m.UpdateCostFunc(hotelName, number, newCost)
}

var _ interfaces.HotelService = (*MockHotelService)(nil)
var _ interfaces.RoomService = (*MockRoomService)(nil)
