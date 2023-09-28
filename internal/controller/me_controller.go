package controller

import (
	"net/http"

	queries "account/config/sqlc"

	"github.com/gin-gonic/gin"
)

type MeController struct{}

// GetMe
//
//	@Summary	获取当前登录用户
//	@Accept		json
//	@Produce	json
//
//	@Security	Bearer
//
//	@Success	200	{object}	api.GetMeResponse
//	@Failure	401	{string}	string	无效的JWT
//	@Router		/api/v1/me [get]
func (ctrl *MeController) Get(c *gin.Context) {
	me, _ := c.Get("me")
	// 断言
	if user, ok := me.(queries.User); !ok {
		c.Status(http.StatusUnauthorized)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"resource": user,
		})
	}

}

func (ctrl *MeController) Create(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *MeController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *MeController) Find(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *MeController) Destory(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *MeController) GetPaged(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *MeController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("v1")
	// 注册路由
	v1.GET("me", ctrl.Get)
}
