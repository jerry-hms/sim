package im

import (
	"github.com/gin-gonic/gin"
	imContro "sim/app/gateway/http/controller/api/v1/im"
	"sim/app/gateway/http/validator/core/data_transfer"
	"sim/app/util/response"
)

type SessionList struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"page_size" form:"page_size"`
}

func (s *SessionList) CheckParams(c *gin.Context) {
	err := c.ShouldBind(&s)
	if err != nil {
		response.ValidatorError(c, s, err)
		return
	}

	if context := data_transfer.DataAddContext(s, "", c); context != nil {
		(&imContro.ImControl{}).SessionList(c)
	} else {
		response.Fail(c, "参数绑定失败", nil)
		return
	}
}
