package dto

type CreateHotelRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateHotelRequest struct {
	Email string `json:"email"`
}

type UpdateRoomRequest struct {
	Number    int     `json:"number"`
	NewCost   float32 `json:"new_cost"`
	UserEmail string  `json:"user_email"`
}
