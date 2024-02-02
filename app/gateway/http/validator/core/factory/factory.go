package factory

import (
	"github.com/gin-gonic/gin"
	"sim/app/core/container"
	"sim/app/gateway/http/validator/core/interf"
)

// Create 创建验证器工厂
func Create(validate string) func(c *gin.Context) {
	valid := container.CreateContainersFactory().Get(validate)
	// 判断是ValidatorInterface接口则返回验证闭包
	if validatorInterface, ok := valid.(interf.ValidatorInterface); ok {
		return validatorInterface.CheckParams
	}
	return nil
}
