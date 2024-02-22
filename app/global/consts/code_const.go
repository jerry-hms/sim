package consts

const (

	// api状态码相关
	RequestSuccess int = 200  // 请求成功状态码
	RequestFail    int = -400 // 请求失败

	// token 相关
	TokenParseError   int = -400100 // token解析错误
	TokenInvalidError int = -400101 // token无效
	TokenExpiredError int = -400102 // token已过期
	TokenFormatError  int = -400103 // token格式错误

	// 验证器相关
	ValidatorCheckParamsFail int = -400300

	// 路由相关
	RouterNotFoundFail int = -400400

	// 系统错误
	SystemError int = -500100
)
