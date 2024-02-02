package im

import (
	"github.com/gin-gonic/gin"
	v1 "sim/app/gateway/http/controller/api/v1"
	"sim/app/gateway/http/validator/core/data_transfer"
	"sim/app/gateway/http/validator/core/verify_params"
	"sim/app/util/response"
)

type Message struct {
	RecvId  uint64 `json:"recv_id" form:"recv_id" validate:"required" valid_msg:"接收人必填"`
	Type    string `json:"type" form:"type" validate:"required" valid_msg:"消息类型必填"`
	Content string `json:"content" form:"content" validate:"required" valid_msg:"发送内容必填"`
	Scene   string `json:"scene" form:"scene" validate:"required" valid_msg:"发送场景必填"`
}

func (m *Message) CheckParams(c *gin.Context) {
	err := c.ShouldBind(&m)
	if err != nil {
		response.ValidatorFail(c, "参数绑定失败")
		return
	}

	verify, err := verify_params.Verify(m)
	if err != nil {
		response.ValidatorFail(c, "验证器故障")
		return
	}
	if verify != "" {
		response.ValidatorFail(c, verify)
		return
	}
	if context := data_transfer.DataAddContext(m, "", c); context != nil {
		(&v1.Chat{}).Send(c)
	} else {
		response.ValidatorFail(c, "参数绑定失败")
		return
	}
}
