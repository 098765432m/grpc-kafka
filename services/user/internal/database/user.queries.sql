-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :exec
INSERT INTO users (
    username, 
    password,
    email,
    phone_number,
    full_name,
    role,
    hotel_id 
) VALUES (
    @username::text, 
    @password::text,
    @email::text,
    @phone_number::text,
    @full_name::text,
    @role::role_enum,
    @hotel_id::text
)

-- name: UpdateUser :one
UPDATE users
SET
    username = @username::text,
    password = @password::text,
    email = @email::text,
    phone_number = @phone_number::text,
    full_name = @full_name::text,
    role = @role::role_enum,
    hotel_id = @hotel_id::text
WHERE id = @id::uuid
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;