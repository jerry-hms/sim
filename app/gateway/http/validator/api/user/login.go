package user

import (
	"github.com/gin-gonic/gin"
	userV1 "sim/app/gateway/http/controller/api/v1/user"
	"sim/app/gateway/http/validator/core/data_transfer"
	"sim/app/util/response"
)

type Login struct {
	BaseField
}

func (l Login) CheckParams(c *gin.Context) {
	err := c.ShouldBind(&l)
	if err != nil {
		response.ValidatorError(c, l, err)
		return
	}

	context := data_transfer.DataAddContext(l, "", c)
	if context == nil {
		response.Fail(c, "数据绑定失败", nil)
	} else {
		(&userV1.User{}).Login(c)
	}
}
