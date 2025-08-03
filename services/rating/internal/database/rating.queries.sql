-- name: GetRatingsByHotel :many
SELECT *
FROM ratings
WHERE hotel_id = $1;

-- name: CreateRating :exec
INSERT INTO ratings
(
    rating,
    hotel_id,
    user_id,
    comment
)
VALUES
(
    @rating::int,
    @hotel_id::text,
    @user_id::text,
    @comment::text
);

-- name: UpdateRating :exec
UPDATE ratings
SET 
    rating = @rating::int,
    hotel_id = @hotel_id::text,
    user_id = @user_id::text,
    comment = @comment::text
WHERE id = @id::text;

-- name: DeleteRating :exec
DELETE FROM ratings
WHERE id = $1;