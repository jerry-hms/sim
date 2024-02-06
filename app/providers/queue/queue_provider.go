package queue

import (
	"sim/app/core/container"
	"sim/app/global/consts"
	"sim/app/listeners"
)

func RegisterQueue() {
	var key string
	contain := container.CreateContainersFactory()
	key = consts.QueuePrefix + "send_message"
	contain.Set(key, &listeners.SendMessageListener{})
}
