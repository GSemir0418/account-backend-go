-- name: CreateItem :one
INSERT INTO items (
  user_id,
  amount,
  kind,
  happened_at,
  tag_ids
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: ListItems :many
SELECT * FROM items
ORDER BY happened_at DESC
OFFSET $1
LIMIT $2;

-- name: CountItems :one
SELECT count(*) FROM items;

-- name: DeleteAllItems :exec
DELETE FROM items;
