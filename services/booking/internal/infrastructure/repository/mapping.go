package booking_repo_mapping

import (
	booking_domain "github.com/098765432m/grpc-kafka/booking/internal/domain"
	booking_repo "github.com/098765432m/grpc-kafka/booking/internal/infrastructure/repository/sqlc/booking"
)

func FromBookingRepoToBookingDomain(bookingRepo booking_repo.Booking) booking_domain.Booking {

	return booking_domain.Booking{
		Id:         bookingRepo.ID.String(),
		CheckIn:    bookingRepo.CheckIn.Time,
		CheckOut:   bookingRepo.CheckOut.Time,
		Total:      int(bookingRepo.Total),
		Status:     string(bookingRepo.Status),
		HotelId:    bookingRepo.HotelID.String(),
		RoomTypeId: bookingRepo.RoomTypeID.String(),
		RoomId:     bookingRepo.RoomID.String(),
		UserId:     bookingRepo.UserID.String(),
		CreatedAt:  bookingRepo.CreateAt.Time,
		UpdatedAt:  bookingRepo.UpdatedAt.Time,
	}
}

func FromBookingsRepoToBookingsDomain(bookingsRepo []booking_repo.Booking) []booking_domain.Booking {

	bookings := make([]booking_domain.Booking, 0, len(bookingsRepo))

	for _, b := range bookingsRepo {
		bookings = append(bookings, FromBookingRepoToBookingDomain(b))
	}

	return bookings
}
