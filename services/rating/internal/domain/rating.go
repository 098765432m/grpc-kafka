package rating_domain

type Rating struct {
	Id      string `json:"id"`
	Score   int    `json:"score"`
	UserId  string `json:"user_id"`
	HotelId string `json:"hotel_id"`
	Comment string `json:"comment,omitempty"`
}
