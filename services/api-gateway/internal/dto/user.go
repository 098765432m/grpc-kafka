package api_dto

type UserImage struct {
	Id       string `json:"id"`
	PublicId string `json:"public_id"`
	Format   string `json:"format"`
}

type UserResponse struct {
	Id          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	FullName    string     `json:"full_name"`
	Role        string     `json:"role,omitempty"`
	HotelId     string     `json:"hotelId,omitempty"`
	Image       *UserImage `json:"image,omitempty"`
}
