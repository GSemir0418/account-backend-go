package controller_test

import (
	queries "account/config/sqlc"
	"account/internal/database"
	"account/internal/router"
	"context"
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

var (
	r *gin.Engine
	q *queries.Queries
	c context.Context
)

func setUpTest(t *testing.T) func(t *testing.T) {
	r = router.New()
	q = database.NewQuery()
	c = context.Background()
	// 删除 User 表
	if err := q.DeleteAllUsers(c); err != nil {
		t.Fatal(err)
	}
	return func(t *testing.T) {
		database.Close()
	}

}

func TestSession(t *testing.T) {
	// 初始化测试环境
	teardownTest := setUpTest(t)
	defer teardownTest(t)
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
	var resBody struct {
		JWT    string `json:"jwt"`
		UserID int32  `json:"userId"`
	}
	fmt.Println("JWT==========")
	fmt.Println(w.Body.String())
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("jwt is not a string")
	}
	// 断言
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, user.ID, resBody.UserID)
}
