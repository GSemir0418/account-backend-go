package controller

import (
	queries "account/config/sqlc"
	"account/internal/database"
	"net/http"
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
//	@Param		amount		body		int				true	"金额（单位：分）"
//	@Param		kind		body		queries.Kind	true	"类型"
//	@Param		happened_at	body		string			true	"发生时间"
//	@Param		tag_ids		body		[]string		true	"标签ID列表"
//	@Success	200			{object}	api.CreateItemResponse
//	@Failure	401			{string}	string	无效的JWT
//	@Failure	422			{string}	string	参数错误
//	@Router		/api/v1/items [post]
func (ctrl *ItemController) Create(c *gin.Context) {
	var body struct {
		Amount     int32        `json:"amount" binding:"required"`
		Kind       queries.Kind `json:"kind" binding:"required"`
		HappenedAt time.Time    `json:"happened_at" binding:"required"`
		TagIds     []int32      `json:"tag_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(422, "参数错误")
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

func (ctrl *ItemController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("v1")
	// 注册路由
	v1.POST("/items", ctrl.Create)
}
