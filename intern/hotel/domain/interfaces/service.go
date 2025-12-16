package interfaces

type HotelService interface {
	GetAll() ([]string, error)
	GetEmail(name string) (string, error)
	Create(name string, email string) error
}

type RoomService interface {
	GetAll(name string) ([]int, error)
	Create(hotelName string, number int, cost float32) error
	GetCost(name string, number int) (float32, error)
	UpdateCost(hotelName string, number int, newCost float32) error
}
