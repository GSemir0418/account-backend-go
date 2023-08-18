package router

import (
	"account/internal/controller"

	"github.com/gin-gonic/gin"
)

// 不要担心 New 是关键字，在go中没有构造函数，所有函数都是平等的
func New() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", controller.Ping)
	return r
}
