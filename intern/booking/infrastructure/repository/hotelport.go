package repository

type HotelPort interface {
	CheckAccuracy(hotel string, roomNumber uint64) (bool, error)
	GetAllRoomsInHotel(hotel string) ([]uint64, error)
}
