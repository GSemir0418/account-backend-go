package api

import (
	queries "account/config/sqlc"
)

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Kind string `json:"kind" binding:"required"`
	Sign string `json:"sign" binding:"required"`
}
type UpdateTagRequest struct {
	Name string `json:"name"`
	Sign string `json:"sign"`
	Kind string `json:"kind"`
}

type CreateTagResponse struct {
	Resource queries.Tag `json:"resource"`
}
type UpdateTagResponse CreateTagResponse
