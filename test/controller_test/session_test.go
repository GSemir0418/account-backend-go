package controller_test

import (
	"account/internal/router"
	"encoding/json"
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
	// 响应记录器
	w := httptest.NewRecorder()
	// json 请求体
	reqBody := gin.H{
		"email": "1@qq.com",
		"code":  "123456",
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
}
