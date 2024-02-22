package im

import (
	"github.com/gin-gonic/gin"
	imControl "sim/app/gateway/http/controller/api/v1/im"
	"sim/app/gateway/http/validator/core/data_transfer"
	"sim/app/util/response"
)

type Message struct {
	RecvId  uint64 `json:"recv_id" form:"recv_id" validate:"required" valid_msg:"接收人必填"`
	Type    string `json:"type" form:"type" validate:"required" valid_msg:"消息类型必填"`
	Content string `json:"content" form:"content" validate:"required" valid_msg:"发送内容必填"`
	Scene   string `json:"scene" form:"scene" validate:"required" valid_msg:"发送场景必填"`
	Url     string `json:"url" form:"url"`
	Height  int64  `json:"height" form:"height"`
	Width   int64  `json:"width" form:"width"`
}

func (m *Message) CheckParams(c *gin.Context) {
	err := c.ShouldBind(&m)
	if err != nil {
		response.ValidatorError(c, m, err)
		return
	}

	if context := data_transfer.DataAddContext(m, "", c); context != nil {
		(&imControl.ImControl{}).Send(c)
	} else {
		response.Fail(c, "参数绑定失败", nil)
		return
	}
}
