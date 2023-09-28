package controller

import (
	"account/api"
	queries "account/config/sqlc"
	"account/internal/database"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ItemController struct{}

func (ctrl *ItemController) Get(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// CreateItem
//
//	@Summary	创建账目
//	@Accept		json
//	@Produce	json
//
//	@Security	Bearer
//
//	@Param		amount		body		int				true	"金额（单位：分）"	example(100)
//	@Param		kind		body		queries.Kind	true	"类型"		example(expenses)
//	@Param		happened_at	body		string			true	"发生时间"		example(2023-09-26T00:00:00Z)
//	@Param		tag_ids		body		[]int		true	"标签ID列表"	example([1,2,3])
//
//	@Success	200			{object}	api.CreateItemResponse
//	@Failure	401			{string}	string	无效的JWT
//	@Failure	422			{string}	string	参数错误
//	@Router		/api/v1/items [post]
func (ctrl *ItemController) Create(c *gin.Context) {
	var body api.CreateItemRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(422, "参数错误", err)
		return
	}
	// 获取当前登录用户
	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	item, err := q.CreateItem(c, queries.CreateItemParams{
		UserID:     user.ID,
		Amount:     body.Amount,
		Kind:       body.Kind,
		HappenedAt: body.HappenedAt,
		TagIds:     body.TagIds,
	})
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resource": item,
	})

}

func (ctrl *ItemController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ItemController) Find(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ItemController) Destory(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ItemController) GetPaged(c *gin.Context) {
	var params api.GetPagedItemsRequest
	pageStr, _ := c.Params.Get("page")
	if page, err := strconv.Atoi(pageStr); err == nil {
		params.Page = int32(page)
	}
	// page 默认值 1
	if params.Page == 0 {
		params.Page = 1
	}
	pageSizeStr, _ := c.Params.Get("page_size")
	if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
		params.PageSize = int32(pageSize)
	}
	// pageSize 默认值 10
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	happenedBefore, has := c.Params.Get("happened_before")
	// 时间参数处理
	if has {
		if t, err := time.Parse(time.RFC3339, happenedBefore); err == nil {
			params.HappenedBefore = t
		}
	}
	happenedAfter, has := c.Params.Get("happened_after")
	if has {
		if t, err := time.Parse(time.RFC3339, happenedAfter); err == nil {
			params.HappenedAfter = t
		}
	}

	q := database.NewQuery()
	items, err := q.ListItems(c, queries.ListItemsParams{
		Offset: (params.Page - 1) * params.Page,
		Limit:  params.PageSize,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	count, err := q.CountItems(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	c.JSON(http.StatusOK, api.GetPagesItemsResponse{
		Resources: items,
		Pager: api.Pager{
			Page:     params.Page,
			PageSize: params.PageSize,
			Total:    count,
		},
	})
}

func (ctrl *ItemController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("v1")
	// 注册路由
	v1.POST("/items", ctrl.Create)
	v1.GET("/items", ctrl.GetPaged)
}
