package controller

import (
	"account/internal/database"
	"account/internal/jwt_helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type MeController struct{}

// GetMe
//
//	@Summary	获取当前登录用户
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	api.GetMeResponse
//	@Failure	401 {string}  无效的JWT
//	@Router		/api/v1/me [get]
func (ctrl *MeController) Get(c *gin.Context) {
	// 从 Header 拿到 jwt
	auth := c.GetHeader("Authorization")
	if len(auth) < 8 {
		c.String(401, "无效的jwt")
		return
	}
	jwtString := auth[7:]
	// 解析 jwtString
	t, err := jwt_helper.ParseJWT(jwtString)
	if err != nil {
		c.String(401, "无效的jwt")
		return
	}
	// 解析结果为 jwt.Token 类型，其内部有 Claims 属性
	// Claims 是JWT中的声明部分
	// 将其断言为 jwt.MapClaims 类型，方便访问jwt声明中的键值对信息
	// 断言语法会返回 值和ok，根据ok判断是否断言成功
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		c.String(401, "无效的jwt")
		return
	}
	// 从声明中获取 user_id
	// userID, ok := claims["user_id"].(int32)
	// 通过断点调试，发现 userID是 float64 类型，需要将其转为 int32 类型（作为query的参数）
	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.String(401, "无效的jwt")
		return
	}
	// 先首先断言为 string，再转为 int 再转为 int32
	// userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.String(401, "无效的jwt")
		return
	}
	q := database.NewQuery()
	u, err := q.FindUser(c, int32(userID))
	if err != nil {
		c.String(401, "无效的jwt")
		return
	}
	c.JSON(200, gin.H{
		"resource": u,
	})

}

func (ctrl *MeController) Create(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *MeController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *MeController) Find(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *MeController) Destory(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *MeController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("v1")
	// 注册路由
	v1.GET("me", ctrl.Get)
}
