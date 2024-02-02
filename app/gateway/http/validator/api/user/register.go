package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	apiV1 "sim/app/gateway/http/controller/api/v1"
	"sim/app/gateway/http/validator/core/data_transfer"
	"sim/app/util/response"
	userPb "sim/idl/pb/user"
)

type Register struct {
	BaseField
	Nickname string `json:"nickname" form:"nickname"`
	Mobile   string `json:"mobile" form:"mobile" binding:"required"`
	Avatar   string `json:"avatar" form:"avatar"`
}

func (r Register) CheckParams(c *gin.Context) {
	if err := c.ShouldBind(&r); err != nil {
		response.ValidatorFail(c, "参数校验失败")
		return
	}
	// 将参数绑定到context上，在后续方法中方便使用
	addBindDataContext := data_transfer.DataAddContext(r, "", c)
	if addBindDataContext == nil {
		response.Fail(c, "参数绑定失败", nil)
	} else {
		// 调用控制器
		var userReq userPb.UserRequest
		jsonStr, _ := json.Marshal(r)
		err := json.Unmarshal(jsonStr, &userReq)
		if err != nil {
			response.Fail(c, "参数转换失败", nil)
		}

		(&apiV1.User{}).Register(c, &userReq)
	}
}
