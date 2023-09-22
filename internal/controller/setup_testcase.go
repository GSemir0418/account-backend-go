package controller

import (
	viper_config "account/config"
	queries "account/config/sqlc"
	"account/internal/database"
	"account/internal/middleware"
	"context"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	r *gin.Engine
	q *queries.Queries
	c context.Context
)

func setUpTestCase(t *testing.T) func(t *testing.T) {
	// 读取 viper 配置
	viper_config.LoadViperConfig()
	// 连接数据库
	database.Connect()
	q = database.NewQuery()
	// 初始化 gin 服务器
	r = gin.Default()
	// 应用中间件
	r.Use(middleware.Me([]string{"/swagger", "/api/v1/session", "/api/v1/validation_codes", "/ping"}))
	// 默认上下文
	c = context.Background()
	// 删除 User 表
	if err := q.DeleteAllUsers(c); err != nil {
		t.Fatal(err)
	}
	// 返回清理函数，开发者自行选择执行
	return func(t *testing.T) {
		database.Close()
	}

}
