-- name: GetHotelById :one
SELECT * 
FROM hotels 
WHERE id = $1;

-- name: CreateHotel :exec
INSERT INTO hotels (name, address) VALUES (@name::text, @address::text);

-- name: GetAll :many
SELECT * 
FROM hotels
LIMIT 6;

-- name: UpdateHotelById :exec
UPDATE hotels
SET 
    name = @name::text,
    address = @address::text
WHERE id = @id::uuid;

-- name: DeleteHotelById :exec
DELETE FROM hotels WHERE id = $1;