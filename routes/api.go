package routes

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"sim/app/gateway/http/controller/websocket"
	"sim/app/gateway/http/middleware/authorization"
	ValidatorFactory "sim/app/gateway/http/validator/core/factory"
	"sim/app/global/consts"
	"sim/app/global/variable"
	"sim/app/util/gin_release"
	"sim/app/util/response"
)

func RegisterRoute() *gin.Engine {
	router := gin.Default()

	// 非调试模式（生产模式） 日志写到日志文件
	if variable.ConfigYml.GetBool("appDebug") == false {
		//1.gin自行记录接口访问日志，不需要nginx，如果开启以下3行，那么请屏蔽第 34 行代码
		//gin.DisableConsoleColor()
		//f, _ := os.Create(variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"))
		//gin.DefaultWriter = io.MultiWriter(f)

		//【生产模式】
		// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
		// 如果部署到生产环境，请使用以下模式：
		// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
		// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
		// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志
		router = gin_release.ReleaseRouter()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}

	router.NoRoute(func(c *gin.Context) {
		response.RouterNotFoundFail(c)
	})

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
		im := api.Group("/im").Use(authorization.CheckLogin)
		{
			// 客户端绑定到ws服务
			im.POST("bind-ws", ValidatorFactory.Create(consts.ValidatorPrefix+"ChatBind"))
			// 消息发送接口
			im.POST("send", ValidatorFactory.Create(consts.ValidatorPrefix+"Send"))
			// 会话列表接口
			im.GET("session-list", ValidatorFactory.Create(consts.ValidatorPrefix+"SessionList"))
		}
	}

	return router
}
