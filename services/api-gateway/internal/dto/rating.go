package api_dto

type RatingImageResponse struct {
	ImageId  string `json:"image_id"`
	PublicId string `json:"public_id"`
	Format   string `json:"format"`
	UserId   string `json:"user_id"`
}

type RatingUserResponse struct {
	UserId   string              `json:"user_id"`
	Username string              `json:"useranme"`
	Image    RatingImageResponse `json:"image,omitempty"`
}

type RatingResponse struct {
	Id      string             `json:"id"`
	Rating  int                `json:"rating"`
	HotelId string             `json:"hotel_id"`
	User    RatingUserResponse `json:"user"`
	Comment string             `json:"comment"`
}
