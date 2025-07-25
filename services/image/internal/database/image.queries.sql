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
)

-- name: GetHotelImages :many
SELECT 
    id,
    public_id,
    format,
    hotel_id
FROM images
WHERE hotel_id = $1;

-- name: DeleteImage :exec
DELETE FROM images WHERE id = $1;