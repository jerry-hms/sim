package routes

import (
	"github.com/gin-gonic/gin"
	"sim/app/gateway/http/controller/websocket"
	"sim/app/gateway/http/middleware/authorization"
	ValidatorFactory "sim/app/gateway/http/validator/core/factory"
	"sim/app/global/consts"
)

func RegisterRoute() *gin.Engine {
	router := gin.Default()

	router.GET("ws", (&websocket.Ws{}).Handle)

	api := router.Group("api/v1")
	{
		// 用户相关接口（不需要验证登录的接口）
		userNeedLogin := api.Group("/user")
		{
			userNeedLogin.POST("register", ValidatorFactory.Create(consts.ValidatorPrefix+"UserRegister"))
			userNeedLogin.POST("login", ValidatorFactory.Create(consts.ValidatorPrefix+"UserLogin"))
		}

		// im相关接口
		chat := api.Group("/im").Use(authorization.CheckLogin)
		{
			// 客户端绑定到ws服务
			chat.POST("bind-ws", ValidatorFactory.Create(consts.ValidatorPrefix+"ChatBind"))
			// 消息发送接口
			chat.POST("send", ValidatorFactory.Create(consts.ValidatorPrefix+"Send"))
		}
	}

	return router
}
