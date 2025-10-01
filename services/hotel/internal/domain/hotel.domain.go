package hotel_domain

type Hotel struct {
	Id      string
	Name    string
	Address string
}

type RoomType struct {
	Id      string
	Name    string
	Price   int
	HotelId string
}

type Room struct {
	Id         string
	Name       string
	Status     string
	RoomTypeId string
	HotelId    string
}
