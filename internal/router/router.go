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

// 加载Controllers的方法，返回 controller 包的 Controller 接口的切片
func loadControllers() []controller.Controller {
	return []controller.Controller{
		&controller.SessionController{},
	}
}

func New() *gin.Engine {
	// 读取 viper 配置
	viper_config.LoadViperConfig()
	// 连接数据库
	database.Connect()
	// 创建路由
	r := gin.Default()
	r.GET("/ping", controller.Ping)

	// 注册路由
	rg := r.Group("/api")
	for _, ctrl := range loadControllers() {
		ctrl.RegisterRoutes(rg)
	}
	// 文档路由及配置
	docs.SwaggerInfo.Version = "1.0"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
