-- name: GetHotelById :one
SELECT * FROM hotels WHERE id = $1;

-- name: CreateHotel :exec
INSERT INTO hotels (name, address) VALUES (@name::text, @address::text);

-- name: GetAll :many
SELECT * FROM hotels;