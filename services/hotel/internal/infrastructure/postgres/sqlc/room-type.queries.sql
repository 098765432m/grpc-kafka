-- name: GetRoomTypesByHotelId :many
SELECT 
    rt.id,
    rt.name,
    rt.price,
    rt.hotel_id,
    COUNT(r.id) AS number_of_rooms 
FROM room_types rt LEFT JOIN rooms r ON rt.id = r.room_type_id
WHERE rt.hotel_id = $1
GROUP BY rt.id
ORDER BY rt.name
LIMIT 10;

-- name: GetRoomTypeById :one
SELECT *
FROM room_types
WHERE id = $1;

-- name: CreateRoomType :exec
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

-- name: DeleteRoomTypeById :exec
DELETE FROM room_types
WHERE id = $1;