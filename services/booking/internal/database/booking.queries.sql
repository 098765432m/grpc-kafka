-- name: GetBookingById :one
SELECT * FROM bookings WHERE id = $1;

-- name: CreateBooking :exec
INSERT INTO bookings (
    check_in,
    check_out,
    total,
    status
) VALUES (
    @check_in::date,
    @check_out::date,
    @total::int,
    @status::BOOKING_STATUS
);

-- name: DeleteBookingById :exec
DELETE FROM bookings WHERE id = $1;