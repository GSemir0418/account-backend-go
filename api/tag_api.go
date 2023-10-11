package api

import (
	queries "account/config/sqlc"
)

type CreateTagRequest struct {
	Name string       `json:"name" binding:"required"`
	Kind queries.Kind `json:"kind" binding:"required"`
	Sign string       `json:"sign" binding:"required"`
}

type CreateTagResponse struct {
	Resource queries.Tag `json:"resource"`
}
