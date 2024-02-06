package listeners

import (
	"fmt"
	"sim/app/services/im/core"
	"sim/app/services/websocket/client"
	"sim/app/util/tools"
)

type SendMessageListener struct {
	RecvId  uint64        `json:"recv_id"`
	Message *core.Message `json:"message"`
}

func (s *SendMessageListener) Handle(params interface{}) {
	fmt.Println("ppp", params)
	err := tools.InterfaceToStruct(params, &s)
	if err != nil {
		fmt.Println("绑定失败", err.Error())
	}
	if err = client.GetWsClient().Send(s.RecvId, s.Message); err == nil {
		fmt.Println("发送成功")
	}

}
