package consts

const (
	ValidatorPrefix string = "From_validator_" // 验证器前缀

	// 状态码相关
	RequestSuccess      int = 0  // 请求成功状态码
	RequestFail         int = -1 // 请求失败
	RequestValidateFail int = -2 //表单验证失败状态码

	// token 相关
	TokenParseError   int = -3 // token解析错误
	TokenInvalidError int = -4 // token格式不正确
	TokenExpiredError int = -5 // token已失效

	TokenParseErrorMsg   string = "token解析错误"
	TokenInvalidErrorMsg string = "无效的token"
	TokenExpiredErrorMsg string = "token已过期"

	// 队列相关
	QueuePrefix string = "sim_queue_" //队列前缀
)
