package repository

type HotelPort interface {
	CheckAccuracy(hotel string, roomNumber uint64) (bool, error)
	GetAllRoomsInHotel(hotel string) ([]uint64, error)
	GetRoomPrice(hotel string, roomNumber uint64) (string, error)
	GetOwnerEmail(hotel string) (string, error)
}
