package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateValidationCodes godoc
// @Summary      发送验证码到用户邮箱
// @Description  生成验证码，发送至用户邮箱
// @Tags         ValidationCode
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       /validation_codes [post]
func CreateValidationCode(c *gin.Context) {
	// 拿到 json 请求体，分三步
	// 1. 声明结构体实例
	var body struct {
		Email string `json:"email"`
	}
	// 2. 将结构体实例绑定到上下文
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(400, "参数错误")
		return
	}
	// 3. 此时结构体实例body就是请求体了
	log.Println("--------------------")
	log.Println(body.Email)
	c.String(http.StatusOK, "pong")
}
