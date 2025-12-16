package interfaces

import "net/http"

type HotelHandler interface {
	GetAllHotels(w http.ResponseWriter, r *http.Request)
	CreateHotel(w http.ResponseWriter, r *http.Request)
	GetOneHotel(w http.ResponseWriter, r *http.Request)
	UpdateRoom(w http.ResponseWriter, r *http.Request)
	GetAllRooms(w http.ResponseWriter, r *http.Request)
	GetRoom(w http.ResponseWriter, r *http.Request)
	CreateRoom(w http.ResponseWriter, r *http.Request)
}
