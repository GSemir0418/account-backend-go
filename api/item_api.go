package api

import (
	queries "account/config/sqlc"
	"time"
)

type CreateItemRequest struct {
	Amount     int32     `json:"amount" binding:"required"`
	Kind       string    `json:"kind" binding:"required"`
	HappenedAt time.Time `json:"happened_at" binding:"required"`
	TagIds     []int32   `json:"tag_ids" binding:"required"`
}

type CreateItemResponse struct {
	Resource queries.Item
}

type GetPagedItemsRequest struct {
	Page           int32     `json:"page"`
	PageSize       int32     `json:"page_size"`
	HappenedAfter  time.Time `json:"happened_after"`
	HappenedBefore time.Time `json:"happened_before"`
}

type GetPagedItemsResponse struct {
	Resources []queries.Item
	Pager     Pager
}

type GetBalanceResponse struct {
	Balance  int `json:"balance"`
	Expenses int `json:"expenses"`
	Income   int `json:"income"`
}

type GetSummaryRequest struct {
	HappenedAfter  time.Time `form:"happened_after" binding:"required"`
	HappenedBefore time.Time `form:"happened_before" binding:"required"`
	Kind           string    `form:"kind" binding:"required,oneof=expenses in_come"`
	GroupBy        string    `form:"group_by" binding:"required"`
}
