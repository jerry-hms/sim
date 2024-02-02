package api_register_validator

import (
	"sim/app/core/container"
	imValidator "sim/app/gateway/http/validator/api/im"
	userValidator "sim/app/gateway/http/validator/api/user"
	"sim/app/global/consts"
)

// Handler 注册验证器
func Handler() {
	// 注册user validator
	registerUserValidator()
	// 注册chat validator
	registerChatValidator()
}

// 注册user的验证器
func registerUserValidator() {
	cont := container.CreateContainersFactory()
	var key string
	// 用户注册
	key = consts.ValidatorPrefix + "UserRegister"
	cont.Set(key, &userValidator.Register{})
	// 用户登录
	key = consts.ValidatorPrefix + "UserLogin"
	cont.Set(key, &userValidator.Login{})
}

// 注册chat的验证器
func registerChatValidator() {
	cont := container.CreateContainersFactory()
	var key string
	// 客户端绑定到ws
	key = consts.ValidatorPrefix + "ChatBind"
	cont.Set(key, &imValidator.Bind{})
	// 发送消息
	key = consts.ValidatorPrefix + "Send"
	cont.Set(key, &imValidator.Message{})
}
