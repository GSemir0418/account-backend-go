package controller_test

import (
	queries "account/config/sqlc"
	"account/internal/database"
	"account/internal/router"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	// 初始化路由: 初始化路由，连接数据库，读取 viper 配置
	r := router.New()
	// 创建真实的验证码
	email := "1@qq.com"
	code := "1234"
	q := database.NewQuery()
	c := context.Background()
	if _, err := q.CreateValidationCode(c, queries.CreateValidationCodeParams{
		Email: email,
		Code:  code,
	}); err != nil {
		log.Fatalln(err)
	}
	// 响应记录器
	w := httptest.NewRecorder()
	// json 请求体
	reqBody := gin.H{
		"email": email,
		"code":  code,
	}
	// stringfy
	bytes, _ := json.Marshal(reqBody)
	// 构造请求
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/session",
		strings.NewReader(string(bytes)),
	)
	// 发送请求
	r.ServeHTTP(w, req)
	// 断言
	assert.Equal(t, w.Code, http.StatusOK)
	// 测试 resBody 中有 jwt
	var resBody struct {
		JWT string `json:"jwt"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("jwt is not a string")
	}
}
