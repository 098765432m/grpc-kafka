-- name: GetUsers :many
SELECT * FROM users
ORDER BY hotel_id
LIMIT 20;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: CheckUserExistsById :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE id = $1
);

-- name: CheckUserByUsername :one
SELECT 
    id,
    username,
    password,
    email,
    role
FROM users WHERE username = $1;

-- name: CreateUser :exec
INSERT INTO users (
    username,
    password,
    address,
    email,
    phone_number,
    full_name,
    role,
    hotel_id
) VALUES (
    @username::text, 
    @password::text,
    @address::text,
    @email::text,
    @phone_number::text,
    @full_name::text,
    @role::role_enum,
    @hotel_id::uuid
);

-- name: UpdateUserById :exec
UPDATE users
SET
    username = @username::text,
    password = @password::text,
    address = @address::text,
    email = @email::text,
    phone_number = @phone_number::text,
    full_name = @full_name::text,
    role = @role::role_enum,
    hotel_id = @hotel_id::uuid
WHERE id = @id::uuid;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1;