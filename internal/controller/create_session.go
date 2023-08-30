package controller

import (
	queries "account/config/sqlc"
	"account/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSession(c *gin.Context) {
	var reqBody struct {
		Email string `json:"email" binding:"required"` // 必填项
		Code  string `json:"code" binding:"required"`  // 必填项
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
	// 返回 jwt
	jwt := "xxxx"
	// 正常情况下需要先创建相应体的结构体，再赋值
	// gin 提供了 H 方法，可以省略定义结构体的步骤，直接定义 json
	c.JSON(http.StatusOK, gin.H{
		"jwt": jwt,
	})

}
