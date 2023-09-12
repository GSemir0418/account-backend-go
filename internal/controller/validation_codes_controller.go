package controller

import (
	queries "account/config/sqlc"
	"account/internal/database"
	"account/internal/email"
	"crypto/rand"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidationCodeController struct{}

func (ctrl *ValidationCodeController) Create(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(400, "参数错误")
		return
	}
	q := database.NewQuery()
	code, err := generateDigits()
	if err != nil {
		log.Println("[GenerateDigits fail]", err)
		c.String(500, "生成验证码失败")
		return
	}
	vc, err := q.CreateValidationCode(c, queries.CreateValidationCodeParams{
		Email: body.Email,
		Code:  code,
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

func (ctrl *ValidationCodeController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ValidationCodeController) Find(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ValidationCodeController) Destory(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ValidationCodeController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("validation_codes", ctrl.Create)
}

// 使用内置库 crypto/rand 生成随机四位验证码
func generateDigits() (string, error) {
	len := 4
	// 开辟一个 4 字节的切片
	b := make([]byte, len)
	// 使用 Read 方法填充切片
	// 此时 b 的类型为 []uint8 其中uint8的范围是 0-255
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// 将uint8数字转换为字符
	digits := make([]byte, len)
	for i := range b {
		// 数字转换为字符编码
		// b[i]%10 就可以得到一个 0-9 的数字
		// '0' 对应的编码是 48 所以要加 48 转换为字符编码
		digits[i] = b[i]%10 + 48
	}
	// [49, 50, 51, 52] 转为字符串就是 "1234"
	return string(digits), nil
}
