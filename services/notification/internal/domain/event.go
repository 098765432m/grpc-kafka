package notification_domain

import "time"

type BookingCreatedEvent struct {
	BookingId string    `json:"booking_id"`
	UserEmail string    `json:"user_email"`
	HotelName string    `json:"hotel_name"`
	CheckIn   time.Time `json:"check_in"`
	CheckOut  time.Time `json:"check_out"`
}
