package controller

import (
	queries "account/config/sqlc"
	"account/internal/database"
	"account/internal/email"
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
	var body struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(400, "参数错误")
		return
	}
	q := database.NewQuery()
	vc, err := q.CreateValidationCode(c, queries.CreateValidationCodeParams{
		Email: body.Email,
		Code:  "123456",
	})
	if err != nil {
		// TODO 暂时没有做校验
		c.Status(400)
		return
	}
	err = email.SendValidationCode(vc.Email, vc.Code)
	if err != nil {
		log.Println("[SendValidationCode fail]", err)
		c.String(500, "发送失败")
		return
	}
	c.Status(http.StatusOK)
}
