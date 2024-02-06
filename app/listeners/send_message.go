package listeners

import (
	"sim/app/global/variable"
	"sim/app/services/im/core"
	"sim/app/services/websocket"
	"sim/app/util/tools"
)

type SendMessageListener struct {
	RecvId  uint64        `json:"recv_id"`
	Message *core.Message `json:"message"`
}

func (s *SendMessageListener) Handle(params interface{}) {
	err := tools.InterfaceToStruct(params, &s)
	if err != nil {
		variable.ZapLog.Error("队列[send_message]数据绑定失败")
	}
	if err = (&websocket.Ws{}).Send(s.RecvId, s.Message); err == nil {
		variable.ZapLog.Info("队列[send_message]发送成功")
	}

}
