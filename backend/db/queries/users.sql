-- backend/db/queries/users.sql

-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: BanUser :exec
UPDATE users SET is_banned = TRUE WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET
    bio = $2,
    avatar_url = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
