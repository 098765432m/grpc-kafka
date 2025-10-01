package booking_domain

import "time"

type Booking struct {
	Id         string
	CheckIn    time.Time
	CheckOut   time.Time
	Total      int
	Status     string
	HotelId    string
	RoomTypeId string
	RoomId     string
	UserId     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
