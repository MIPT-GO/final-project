package interfaces

import (
	"final-project/intern/hotel/domain/entity"
)

type HotelRepository interface {
	GetAll() ([]entity.Hotel, error)
	GetByName(name string) (entity.Hotel, error)
	AddNewHotel(entity.Hotel) error
	ExistsByName(name string) (bool, error)
}

type RoomRepository interface {
	GetAll(hotelName string) ([]entity.Room, error)
	GetCost(hotelName string, number int) (float32, error)
	AddNewRoom(entity.Room) error
	UpdateCost(hotelName string, number int, newCost float32) error
}
