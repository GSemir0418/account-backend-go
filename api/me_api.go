package api

import (
	queries "account/config/sqlc"
)

type GetMeResponse struct {
	Resource queries.User
}
