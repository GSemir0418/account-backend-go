-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
OFFSET $1
LIMIT $2;

-- name: CreateUser :one
INSERT INTO users (email) values ($1) RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET email = $2 WHERE id = $1;