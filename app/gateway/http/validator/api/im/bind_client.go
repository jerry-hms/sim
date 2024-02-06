package im

import (
	"github.com/gin-gonic/gin"
	apiV1 "sim/app/gateway/http/controller/api/v1/im"
	"sim/app/gateway/http/validator/core/data_transfer"
	"sim/app/util/response"
)

type Bind struct {
	ClientId string `json:"client_id" form:"client_id" binding:"required" err_msg:"客户端id必填"`
}

func (b Bind) CheckParams(c *gin.Context) {
	// 参数绑定
	if err := c.ShouldBind(&b); err != nil {
		response.ValidatorError(c, b, err)
		return
	}
	//// 参数校验
	//msg, err := verify_params.Verify(b)
	//if err != nil {
	//	fmt.Println("验证参数错误")
	//}
	//if msg != "" {
	//	response.ValidatorFail(c, msg)
	//	return
	//}
	context := data_transfer.DataAddContext(b, "", c)
	if context == nil {
		response.Fail(c, "参数绑定失败", nil)
		return
	} else {
		(&apiV1.ImControl{}).BindToWs(c)
	}
}
