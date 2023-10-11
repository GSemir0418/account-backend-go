-- name: CreateTag :one
INSERT INTO tags (
  user_id,
  sign,
  kind,
  name
) VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;