package api

type Pager struct {
	Page     int32
	PageSize int32
	Total    int64
}

type ErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}
