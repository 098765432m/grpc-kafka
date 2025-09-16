-- name: GetBookingById :one
SELECT * FROM bookings WHERE id = $1;

-- name: GetBookingsByRoomId :many
SELECT * FROM bookings WHERE room_id = $1;

-- name: GetNumberOfOccupiedRooms :many
SELECT 
    room_type_id,
    COUNT(DISTINCT room_id) AS number_of_occupied_rooms
FROM bookings 
WHERE
    room_type_id = ANY(@room_type_ids::uuid[])
    -- AND ( date_trunc('day', @new_check_in::date) < date_trunc('day', check_out) AND date_trunc('day', @new_check_out::date) > date_trunc('day', check_in) );
    -- AND (@new_check_in::date < check_out::date AND @new_check_out::date > check_in::date)
    AND daterange(check_in, check_out, '[]') && daterange(@check_in::date, @check_out::date, '[]')
GROUP BY room_type_id;

-- name: GetNumberOfOccupiedRoomsByHotelIds :many
SELECT
    room_type_id,
    COUNT(DISTINCT room_id) AS number_of_occupied_rooms
FROM bookings 
WHERE
    hotel_id = ANY(@hotel_ids::uuid[])
    -- AND ( date_trunc('day', @new_check_in::date) < date_trunc('day', check_out) AND date_trunc('day', @new_check_out::date) > date_trunc('day', check_in) );
    -- AND (@new_check_in::date < check_out::date AND @new_check_out::date > check_in::date)
    AND daterange(check_in, check_out, '[]') && daterange(@check_in::date, @check_out::date, '[]')
GROUP BY room_type_id;

-- name: GetUnavailableRoomsByRoomTypeId :many
SELECT room_id
FROM bookings
WHERE 
    room_type_id = @room_type_id::uuid
    AND daterange(check_in, check_out, '[]') && daterange(@check_in::date, @check_out::date, '[]');

-- name: CreateBooking :exec
INSERT INTO bookings (
    check_in,
    check_out,
    total,
    status,
    room_type_id,
    user_id,
    room_id
) VALUES (
    @check_in::date,
    @check_out::date,
    @total::int,
    @status::BOOKING_STATUS,
    @room_type_id::uuid,
    @user_id::uuid,
    @room_id::uuid
);

-- name: DeleteBookingById :exec
DELETE FROM bookings WHERE id = $1;

-- name: DeleteBookingsByIds :exec
DELETE FROM bookings WHERE ids = ANY(@booking_ids::uuid[]);