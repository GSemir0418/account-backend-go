package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping
//	@Summary		服务健康度
//	@Description	如果返回 200，说明服务正常运行
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		500
//	@Router			/ping [get]
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
