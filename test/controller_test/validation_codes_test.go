package controller_test

import (
	"account/internal/database"
	"account/internal/router"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateValidationCode(t *testing.T) {
	// 注意要先在 New 中连接数据库，然后再查询
	r := router.New()
	viper.Set("email.smtp.port", "1025")
	viper.Set("email.smtp.host", "localhost")
	q := database.NewQuery()
	c := context.Background()
	count1, _ := q.CountValidationCodes(c, "845217811@qq.com")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/validation_codes", strings.NewReader(`{"email": "845217811@qq.com"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	count2, _ := q.CountValidationCodes(c, "845217811@qq.com")

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, count1, count2-1)
}
