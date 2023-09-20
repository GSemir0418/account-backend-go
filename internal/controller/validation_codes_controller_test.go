package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateValidationCode(t *testing.T) {
	// 初始化测试环境
	teardownTest := setUpTestCase(t)
	defer teardownTest(t)
	// 注册路由
	vcc := ValidationCodeController{}
	vcc.RegisterRoutes(r.Group("/api"))

	viper.Set("email.smtp.port", "1025")
	viper.Set("email.smtp.host", "localhost")
	count1, _ := q.CountValidationCodes(c, "845217811@qq.com")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/validation_codes", strings.NewReader(`{"email": "845217811@qq.com"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	count2, _ := q.CountValidationCodes(c, "845217811@qq.com")

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, count1, count2-1)
}
