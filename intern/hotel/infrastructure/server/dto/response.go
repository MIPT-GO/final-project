package dto

type AllHotelsResponse struct {
	Hotels []string `json:"hotels"`
}

type AllRoomsResponse struct {
	Numbers []int `json:"numbers"`
}

type GetHotelResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
