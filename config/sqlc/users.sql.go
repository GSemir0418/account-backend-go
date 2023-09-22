// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: users.sql

package queries

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (email) values ($1) RETURNING id, email, created_at, updated_at
`

func (q *Queries) CreateUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const findUser = `-- name: FindUser :one
SELECT id, email, created_at, updated_at FROM users
WHERE id = $1
`

func (q *Queries) FindUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, findUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, email, created_at, updated_at FROM users
WHERE email = $1
`

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, created_at, updated_at FROM users
ORDER BY id
OFFSET $1
LIMIT $2
`

type ListUsersParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users SET email = $2 WHERE id = $1
`

type UpdateUserParams struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser, arg.ID, arg.Email)
	return err
}
