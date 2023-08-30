package router

import (
	viper_config "account/config"
	"account/internal/controller"
	"account/internal/database"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	// 也可以手动添加文档信息
	"account/docs"
)

func New() *gin.Engine {
	// 读取 viper 配置
	viper_config.LoadViperConfig()
	// 连接数据库
	database.Connect()
	// 创建路由
	r := gin.Default()
	r.GET("/api/v1/ping", controller.Ping)
	r.POST("/api/v1/validation_codes", controller.CreateValidationCode)
	r.POST("/api/v1/session", controller.CreateSession)
	// 文档路由及配置
	docs.SwaggerInfo.Version = "1.0"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
