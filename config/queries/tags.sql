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
  user_id = @user_id,
  name = CASE WHEN @name::varchar = '' THEN name ELSE @name END,
  sign = CASE WHEN @sign::varchar = '' THEN sign ELSE @sign END,
  kind = CASE WHEN @kind::varchar = '' THEN kind ELSE @kind END
WHERE id = @id
RETURNING id, user_id, name, sign, kind, deleted_at, created_at, updated_at;

-- name: DeleteTag :exec
UPDATE tags 
SET 
  deleted_at = now()
WHERE id = @id;

-- name: ListTags :many
SELECT * FROM tags
WHERE kind = @kind AND deleted_at IS NULL AND user_id = @user_id
ORDER BY created_at
OFFSET $1
LIMIT $2;

-- name: CountTags :one
SELECT count(*) FROM tags
WHERE deleted_at IS NULL;

-- name: DeleteAllTags :exec
DELETE FROM tags;