package controller

import (
	"account/api"
	queries "account/config/sqlc"
	"account/internal/database"

	"github.com/gin-gonic/gin"
)

type TagController struct {
}

func (ctrl *TagController) Get(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// CreateTag
//
//	@Summary	创建标签
//	@Accept		json
//	@Produce	json
//
//	@Security	Bearer
//
//	@Param		name		body		string				true	"金额（单位：分）"	example(通勤)
//	@Param		kind		body		queries.Kind	true	"类型"		example(expenses)
//	@Param		sign	body		string			true	"符号"		example(😈)
//
//	@Success	200			{object}	api.CreateTagResponse
//	@Failure	401			{string}	string	无效的JWT
//	@Failure	422			{string}	string	参数错误
//	@Router		/api/v1/tags [post]
func (ctrl *TagController) Create(c *gin.Context) {
	var reqBody api.CreateTagRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.String(422, "参数错误")
		return
	}
	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	tag, err := q.CreateTag(c, queries.CreateTagParams{
		Kind:   reqBody.Kind,
		Name:   reqBody.Name,
		Sign:   reqBody.Sign,
		UserID: user.ID,
	})
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, api.CreateTagResponse{Resource: tag})
}

func (ctrl *TagController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *TagController) Find(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *TagController) Destory(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *TagController) GetPaged(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *TagController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/tags", ctrl.Create)
}
