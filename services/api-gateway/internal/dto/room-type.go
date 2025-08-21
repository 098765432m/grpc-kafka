package api_dto

type RoomTypeImage struct {
	Id         string `json:"id"`
	PublicId   string `json:"public_id"`
	Format     string `json:"format"`
	RoomTypeId string `json:"room_type_id"`
}

type RoomTypeResponse struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	Price   int             `json:"price"`
	HotelId string          `json:"hotel_id"`
	Images  []RoomTypeImage `json:"images"`
}

type GetNumberOfAvailableRoomsDtoResponse struct {
	Id                     string          `json:"id"`
	Name                   string          `json:"name"`
	Price                  int             `json:"price"`
	HotelId                string          `json:"hotel_id"`
	Images                 []RoomTypeImage `json:"images"`
	NumberOfAvailableRooms int8            `json:"number_of_available_rooms"`
}
