package user_domain

// TODO: Maybe create Role Type enum here

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	Role        string `json:"role"`
	HotelId     string `json:"hotel_id,omitempty"`
}
