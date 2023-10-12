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

-- name: UpdateTag :one
UPDATE tags  
SET 
  user_id = $2,
  name = $3,
  sign = $4,
  kind = $5
WHERE id = $1
RETURNING id, user_id, name, sign, kind, deleted_at, created_at, updated_at;