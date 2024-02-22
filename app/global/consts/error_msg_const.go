package consts

const (
	// token相关错误
	TokenParseErrorMsg   string = "请传入有效的token" // token解析错误
	TokenInvalidErrorMsg string = "无效的token"    // token无效
	TokenExpiredErrorMsg string = "token已过期"    // token已过期
	TokenFormatErrorMsg  string = "token格式不正确"  // token格式错误

	// websocket相关错误
	WebsocketOnOpenFailMsg      string = "websocket连接失败"
	WebsocketUpgradeFailMsg     string = "websocket协议升级失败"
	WebsocketSendMessageFailMsg string = "websocket发送消息失败"
	WebsocketReadPumpFailMsg    string = "websocket ReadPump 协程读取消息错误"
	WebsocketClientLogoutMsg    string = "websocket客户端已下线"
	WebsocketHeartFailMsg       string = "websocket心跳包检测协程错误"

	// 验证器相关
	ValidatorCheckParamsFailMsg string = "参数校验失败"

	// 路由相关
	RouterNotFoundFailMsg string = "请求地址不存在"

	// 服务器内部错误
	SystemErrorMsg string = "服务器内部代码错误"

	// 队列相关
	QueueDataTransErrorFormatMsg             string = "队列数据解析错误，请检查数据格式是否正确"
	QueueNotImplementInterfaceErrorFormatMsg string = "%s is not implement QueueInterface"
)
