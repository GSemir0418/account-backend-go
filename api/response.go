package api

import queries "account/config/sqlc"

type GetMeResponse struct {
	Resource queries.User
}
type CreateItemResponse struct {
	Resource queries.User
}
