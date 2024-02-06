package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	userV1 "sim/app/gateway/http/controller/api/v1/user"
	"sim/app/gateway/http/validator/core/data_transfer"
	"sim/app/util/response"
	userPb "sim/idl/pb/user"
)

type Register struct {
	BaseField
	Nickname string `json:"nickname" form:"nickname" err_msg:"请输入用户昵称"`
	Mobile   string `json:"mobile" form:"mobile" binding:"required" err_msg:"请输入手机号"`
	Avatar   string `json:"avatar" form:"avatar"`
}

func (r Register) CheckParams(c *gin.Context) {
	if err := c.ShouldBind(&r); err != nil {
		response.ValidatorError(c, r, err)
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

		(&userV1.User{}).Register(c, &userReq)
	}
}
