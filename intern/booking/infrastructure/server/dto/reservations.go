package dto

type Reservation struct {
	Email  string `json:"email"`
	Start  string `json:"start"`
	End    string `json:"end"`
	Hotel  string `json:"hotel"`
	Number uint64 `json:"number"`
}

type ReservationsResponse struct {
	Reservations []Reservation `json:"reservations"`
}
