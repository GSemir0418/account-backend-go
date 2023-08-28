package controller_test

import (
	"account/internal/router"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateValidationCode(t *testing.T) {
	// 测试环境不会运行 main.go
	// 且运行了也无法正确读取 viper 的配置（路径）
	// 所以将 viper 的配置文件统一存放到 $HOME/.account 项目目录下
	// 使用绝对路径读取配置文件即可
	// 封装Viper的读取逻辑 在 router 中引入
	r := router.New()
	w := httptest.NewRecorder()
	// 请求体数据需要使用 strings.NewReader 构造
	req, _ := http.NewRequest("POST", "/api/v1/validation_codes", strings.NewReader(`{"email": "845217811@qq.com"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// assert.Equal(t, "123456", w.Body.String())
}
