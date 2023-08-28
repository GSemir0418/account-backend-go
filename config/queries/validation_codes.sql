-- name: CreateValidationCode :one
INSERT INTO validation_codes (email, code) values ($1, $2) RETURNING *;