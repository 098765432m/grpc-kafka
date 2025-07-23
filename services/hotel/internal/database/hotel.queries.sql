-- name: GetHotelById :one
SELECT * FROM hotels WHERE id = $1;

-- name: CreateHotel :exec
INSERT INTO hotels (name) VALUES (@name::text);

-- name: GetAll :many
SELECT * FROM hotels;