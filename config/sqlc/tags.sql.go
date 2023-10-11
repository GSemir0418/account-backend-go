// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: tags.sql

package queries

import (
	"context"
)

const createTag = `-- name: CreateTag :one
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
RETURNING id, user_id, name, sign, kind, deleted_at, created_at, updated_at
`

type CreateTagParams struct {
	UserID int32  `json:"user_id"`
	Sign   string `json:"sign"`
	Kind   Kind   `json:"kind"`
	Name   string `json:"name"`
}

func (q *Queries) CreateTag(ctx context.Context, arg CreateTagParams) (Tag, error) {
	row := q.db.QueryRowContext(ctx, createTag,
		arg.UserID,
		arg.Sign,
		arg.Kind,
		arg.Name,
	)
	var i Tag
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Sign,
		&i.Kind,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
