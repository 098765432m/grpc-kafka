-- name: GetRoomTypesByHotelId :many
SELECT *
FROM room_types rt
WHERE hotel_id = $1
ORDER BY rt.name
LIMIT 10;

-- name: GetRoomTypeById :one
SELECT *
FROM room_types
WHERE id = $1;

-- CreateRoomType :exec
INSERT INTO room_types
(
    name,
    price,
    hotel_id
)
VALUES
(
    @name::text,
    @price::int,
    @hotel_id::uuid
);

-- DeleteRoomTypeById :exec
DELETE FROM room_types
WHERE id = $1;