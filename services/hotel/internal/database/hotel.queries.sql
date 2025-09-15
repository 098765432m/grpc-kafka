-- name: GetHotelById :one
SELECT * 
FROM hotels 
WHERE id = $1;

-- name: CreateHotel :exec
INSERT INTO hotels (name, address) VALUES (@name::text, @address::text);

-- name: GetAll :many
SELECT * 
FROM hotels
LIMIT 20;

-- name: GetHotelsByAddress :many
SELECT *
FROM hotels h
WHERE 
    (
        @address::text IS NULL
        OR unaccent(h.address) ILIKE unaccent(@address::text)
    ) 
    AND 
    (
        @hotel_name::text IS NULL
        OR h.name ILIKE @hotel_name::text
    )
LIMIT 20;

-- name: FilterHotels :many
SELECT 
    h.id,
    MIN(rt.price)
FROM hotels h LEFT JOIN room_types rt ON h.id = rt.hotel_id
WHERE 
    rt.id = ANY(@room_type_ids::uuid[])
    AND @number_of_occupied_rooms::int < 
    (
        SELECT COUNT(rt.id)
        FROM room_types rt LEFT JOIN rooms r ON rt.id = r.room_type_id
    )
    AND
    (
        sqlc.narg('min_price')::int IS NULL
        OR sqlc.narg('max_price')::int IS NULL
        OR rt.price BETWEEN @min_price AND @max_price
    )
GROUP BY h.id;

-- name: UpdateHotelById :exec
UPDATE hotels
SET 
    name = @name::text,
    address = @address::text
WHERE id = @id::uuid;

-- name: DeleteHotelById :exec
DELETE FROM hotels WHERE id = $1;