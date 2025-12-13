package entity

type Hotel struct {
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

func NewHotelEntity(name, email string) *Hotel {
	return &Hotel{
		Name:  name,
		Email: email,
	}
}
