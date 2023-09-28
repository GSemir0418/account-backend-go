package controller

import (
	"account/api"
	queries "account/config/sqlc"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	// 初始化测试环境
	teardownTest := setUpTestCase(t)
	defer teardownTest(t)
	// 注册路由
	sc := SessionController{}
	sc.RegisterRoutes(r.Group("/api"))
	// 创建真实的验证码
	email := "1@qq.com"
	code := "1234"
	if _, err := q.CreateValidationCode(c, queries.CreateValidationCodeParams{
		Email: email,
		Code:  code,
	}); err != nil {
		log.Fatalln(err)
	}
	user, err := q.CreateUser(c, email)
	if err != nil {
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
	// 测试 resBody 中有 jwt
	var resBody api.CreateSessionResponse
	fmt.Println("JWT==========")
	fmt.Println(w.Body.String())
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("jwt is not a string")
	}
	// 断言
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, user.ID, resBody.UserID)
}
