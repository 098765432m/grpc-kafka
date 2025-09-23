-- name: GetImageById :one
SELECT * FROM images WHERE id = $1;

-- name: UploadHotelImage :exec
INSERT INTO images (
    public_id,
    format,
    hotel_id
) VALUES (
    @public_id::text,
    @format::text,
    @hotel_id::uuid
);

-- name: UploadUserImage :exec
INSERT INTO images (
    public_id,
    format,
    user_id
) VALUES (
    @public_id::text,
    @format::text,
    @user_id::uuid
);

-- name: UploadRoomTypeImage :exec
INSERT INTO images (
    public_id,
    format,
    room_type_id
) VALUES (
    @public_id::text,
    @format::text,
    @room_type_id::uuid
);

-- name: GetImagesByHotelId :many
SELECT *
FROM images
WHERE hotel_id = $1;

-- name: GetImageByUserId :one
SELECT *
FROM images
WHERE user_id = $1;

-- name: GetImagesByRoomTypeId :many
SELECT *
FROM images
WHERE room_type_id = $1;

-- name: GetImagesByRoomTypeIds :many
SELECT *
FROM images
WHERE room_type_id = ANY(@room_type_ids::uuid[]);

-- name: GetImagesByHotelIds :many
SELECT *
FROM images
WHERE hotel_id = ANY(@hotel_ids::uuid[]);

-- name: GetImagesByUserIds :many
SELECT *
FROM images
WHERE user_id = ANY(@user_ids::uuid[]);

-- name: DeleteImage :one
DELETE FROM images WHERE id = $1
RETURNING public_id;

-- name: DeleteImages :many
DELETE FROM images WHERE id = ANY($1::uuid[])
RETURNING public_id;