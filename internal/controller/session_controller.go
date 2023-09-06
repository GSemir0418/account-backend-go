package controller

import (
	queries "account/config/sqlc"
	"account/internal/database"
	"account/internal/jwt_helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionController struct{}

func (ctrl *SessionController) Create(c *gin.Context) {
	var reqBody struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	// 把JSON格式的请求体 转换为 go 的结构体
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.String(http.StatusBadRequest, "无效的参数")
		return
	}
	// 查询验证码是否存在且有效
	q := database.NewQuery()
	if _, err := q.FindValidationCode(c, queries.FindValidationCodeParams{
		Email: reqBody.Email,
		Code:  reqBody.Code,
	}); err != nil {
		c.String(http.StatusBadRequest, "无效的验证码")
		return
	}
	// 查询 User，有则返回，没有就创建
	user, err := q.FindUserByEmail(c, reqBody.Email)
	if err != nil {
		user, err = q.CreateUser(c, reqBody.Email)
		if err != nil {
			log.Println("Create User error", err)
			c.String(http.StatusInternalServerError, "请稍后再试")
			return
		}
	}
	// 返回 jwt
	jwt, err := jwt_helper.GenerateJWT(int(user.ID))
	if err != nil {
		log.Println("Generate JWT Error", err)
		c.String(http.StatusInternalServerError, "请稍后再试")
		return
	}
	// 正常情况下需要先创建相应体的结构体，再赋值
	// gin 提供了 H 方法，可以省略定义结构体的步骤，直接定义 json
	c.JSON(http.StatusOK, gin.H{
		"jwt":    jwt,
		"userId": user.ID,
	})

}

func (ctrl *SessionController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *SessionController) Find(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *SessionController) Destory(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *SessionController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("session", ctrl.Create)
}
