package dto

type CreateHotelRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateRoomRequest struct {
	Number  int     `json:"number"`
	NewCost float32 `json:"new_cost"`
}

type CreateRoomRequest struct {
	Number int     `json:"number"`
	Cost   float32 `json:"cost"`
}
