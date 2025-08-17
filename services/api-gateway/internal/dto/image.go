package api_dto

type ImageResponse struct {
	Id         string `json:"id"`
	PublicId   string `json:"public_id"`
	Format     string `json:"format"`
	HotelId    string `json:"hotel_id,omitempty"`
	UserId     string `json:"user_id,omitempty"`
	RoomTypeId string `json:"room_type_id,omitempty"`
}
