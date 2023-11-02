package controller

import (
	"account/api"
	queries "account/config/sqlc"
	"account/internal/database"
	"net/http"
	"strconv"

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
//	@Param		kind		body		string	true	"类型"		example(expenses)
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
	var reqBody api.UpdateTagRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.String(422, "参数错误")
		return
	}
	idString, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.String(422, "参数错误")
		return
	}
	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	tag, err := q.UpdateTag(c, queries.UpdateTagParams{
		ID:     int32(id),
		Kind:   reqBody.Kind,
		Name:   reqBody.Name,
		Sign:   reqBody.Sign,
		UserID: user.ID,
	})
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, api.UpdateTagResponse{Resource: tag})
}

func (ctrl *TagController) Find(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *TagController) Destory(c *gin.Context) {
	idStr, bool := c.Params.Get("id")
	if bool == false {
		c.String(422, "参数错误")
		return
	}
	// string 转 int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(422, "参数错误")
		return
	}
	q := database.NewQuery()
	err = q.DeleteTag(c, int32(id))
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (ctrl *TagController) GetPaged(c *gin.Context) {
	me, _ := c.Get("me")
	user, ok := me.(queries.User)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	var params api.GetPagedTagsRequest
	params.Page = 1
	params.PageSize = 10

	pageStr := c.Request.URL.Query().Get("page")
	if page, err := strconv.Atoi(pageStr); err == nil {
		params.Page = int32(page)
	}
	pageSizeStr := c.Request.URL.Query().Get("page_size")
	if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
		params.PageSize = int32(pageSize)
	}
	kind, _ := c.Params.Get("kind")
	if kind == "" {
		kind = "expenses"
	}

	q := database.NewQuery()
	tags, err := q.ListTags(c, queries.ListTagsParams{
		Offset: (params.Page - 1) * params.PageSize,
		Limit:  params.PageSize,
		Kind:   kind,
		UserID: user.ID,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	count, err := q.CountTags(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	c.JSON(http.StatusOK, api.GetPagedTagsResponse{
		Resources: tags,
		Pager: api.Pager{
			Page:     params.Page,
			PageSize: params.PageSize,
			Total:    count,
		},
	})
}

func (ctrl *TagController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/tags", ctrl.Create)
	v1.PATCH("/tags/:id", ctrl.Update)
	v1.DELETE("/tags/:id", ctrl.Destory)
	v1.GET("/tags", ctrl.GetPaged)
}
