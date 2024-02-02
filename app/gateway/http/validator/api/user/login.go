package user

import (
	"github.com/gin-gonic/gin"
	apiV1 "sim/app/gateway/http/controller/api/v1"
	"sim/app/gateway/http/validator/core/data_transfer"
	"sim/app/util/response"
)

type Login struct {
	BaseField
}

func (l Login) CheckParams(c *gin.Context) {
	err := c.ShouldBind(&l)
	if err != nil {
		response.ValidatorFail(c, err.Error())
		return
	}

	context := data_transfer.DataAddContext(l, "", c)
	if context == nil {
		response.ValidatorFail(c, "数据绑定失败")
	} else {
		(&apiV1.User{}).Login(c)
	}
}
