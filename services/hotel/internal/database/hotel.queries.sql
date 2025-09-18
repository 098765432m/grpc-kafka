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
SELECT h.id
FROM hotels h
WHERE 
    (
        @address::text IS NULL
        OR unaccent(h.address) ILIKE unaccent('%' || @address::text || '%')
    ) 
    AND 
    (
        @hotel_name::text IS NULL
        OR unaccent(h.name) ILIKE unaccent( '%' || @hotel_name::text || '%')
    )
LIMIT 20;

-- name: FilterHotels :many
SELECT 
    h.id,
    h.name,
    h.address,
    MIN(rt.price) AS min_price
FROM hotels h LEFT JOIN room_types rt ON h.id = rt.hotel_id AND 
    rt.id = ANY(@room_type_ids::uuid[])
WHERE 
    (
        sqlc.narg('min_price')::int IS NULL
        OR sqlc.narg('max_price')::int IS NULL
        OR rt.price BETWEEN @min_price AND @max_price
    )
GROUP BY h.id
HAVING MIN(rt.price) > 0;

-- name: UpdateHotelById :exec
UPDATE hotels
SET 
    name = @name::text,
    address = @address::text
WHERE id = @id::uuid;

-- name: DeleteHotelById :exec
DELETE FROM hotels WHERE id = $1;