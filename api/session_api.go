package api

type CreateSessionRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type CreateSessionResponse struct {
	Jwt    string `json:"jwt"`
	UserID int32  `json:"user_id"`
}
