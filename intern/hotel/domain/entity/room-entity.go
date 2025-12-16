package entity

type Room struct {
	HotelName string  `json:"hotel_name" db:"hotel_name"`
	Number    int     `json:"number" db:"number"`
	Cost      float32 `json:"cost" db:"cost"`
}

func NewRoomEntity(hotelName string, number int, cost float32) *Room {
	return &Room{
		HotelName: hotelName,
		Number:    number,
		Cost:      cost,
	}
}
