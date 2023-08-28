package router

import (
	"account/internal/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	// 也可以手动添加文档信息
	"account/docs"
)

func New() *gin.Engine {
	r := gin.Default()
	r.GET("/api/v1/ping", controller.Ping)
	r.POST("/api/v1/validation_codes", controller.CreateValidationCode)
	docs.SwaggerInfo.Version = "1.0"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
