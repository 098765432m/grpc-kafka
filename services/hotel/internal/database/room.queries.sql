-- name: GetRoomsByHotelId :many
SELECT *
FROM rooms r
WHERE hotel_id = $1
ORDER BY r.name
LIMIT 20;

-- name: GetRoomsByRoomTypeId :many
SELECT *
FROM rooms r
WHERE room_type_id = $1
ORDER BY r.name
LIMIT 20;

-- name: GetRoomsById :one
SELECT *
FROM rooms
WHERE id = $1;

-- name: CreateRoom :exec
INSERT INTO rooms
(
    name,
    status,
    hotel_id,
    room_type_id
)
VALUES
(
    @name::text,
    @status::text,
    @hotel_id::uuid,
    @room_type_id::uuid
);

-- name: DeleteRoomById :exec
DELETE FROM rooms
WHERE id = $1;