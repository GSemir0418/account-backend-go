package middleware

import (
	"account/internal/database"
	"account/internal/jwt_helper"
	"fmt"
	"strings"

	queries "account/config/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Me(whiteList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 检测白名单
		for _, s := range whiteList {
			if has := strings.HasPrefix(path, s); has {
				c.Next()
				return
			}
		}

		user, err := getMe(c)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		// 将 me 放到上下文中，作为全局变量
		c.Set("me", user)
		c.Next()
	}

}

func getMe(c *gin.Context) (queries.User, error) {
	var user queries.User
	auth := c.GetHeader("Authorization")
	if len(auth) < 8 {
		return user, fmt.Errorf("JWT为空")
	}
	jwtString := auth[7:]
	t, err := jwt_helper.ParseJWT(jwtString)
	if err != nil {
		return user, fmt.Errorf("无效的jwt")
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return user, fmt.Errorf("无效的jwt")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return user, fmt.Errorf("无效的jwt")
	}
	q := database.NewQuery()
	u, err := q.FindUser(c, int32(userID))
	if err != nil {
		return user, fmt.Errorf("无效的jwt")
	}
	return u, nil
}
