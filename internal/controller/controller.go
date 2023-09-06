package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Find(c *gin.Context)
	Destory(c *gin.Context)
	RegisterRoutes(rg *gin.RouterGroup)
}
