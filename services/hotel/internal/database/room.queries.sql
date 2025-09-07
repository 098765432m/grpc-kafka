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

-- name: GetListOfAvailableRoomsByRoomTypeId :many
SELECT id
FROM rooms r
WHERE r.room_type_id = @room_type_id::uuid
    AND r.status = 'AVAILABLE'
ORDER BY r.name
LIMIT @number_of_rooms::int;

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

-- name: ChangeStatusRoomsByIds :exec
UPDATE rooms
SET status = @status::room_status
WHERE id = ANY(@room_ids::uuid[]);