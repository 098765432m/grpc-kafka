-- name: GetImageById :one
SELECT * FROM images WHERE id = $1;

-- name: UploadImage :exec
INSERT INTO images (
    public_id,
    format,
    hotel_id
) VALUES (
    $1,
    $2,
    $3
);

-- name: GetImagesByHotelId :many
SELECT 
*
FROM images
WHERE hotel_id = $1;

-- name: GetImagesByHotelIds :many
SELECT *
FROM images
WHERE hotel_id = ANY(@hotel_ids::uuid[]);

-- name: DeleteImage :exec
DELETE FROM images WHERE id = $1;

-- name: DeleteImages :exec
DELETE FROM images WHERE id = ANY($1::uuid[]);