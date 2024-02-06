package queue

import (
	"encoding/json"
	"fmt"
	"sim/app/core/container"
	"sim/app/global/consts"
	"sim/app/global/variable"
	"sim/app/util/queue/interf"
	"sim/app/util/rabbitmq/publish_subscribe"
)

const QueueMessageFormat string = `{"event": "%s", "data": %s}` // 队列消息传输格式定义

// Transfer 将队列转发到对应的处理方法
func Transfer(receiveData string) {
	var res map[string]interface{}
	event := res["event"].(string)

	if err := json.Unmarshal([]byte(receiveData), &res); err != nil {
		fmt.Println("event数据解析失败", err.Error())
		variable.ZapLog.Error(fmt.Sprintf(""))
	}
	params := res["data"]
	fmt.Println("params", params)
	fmt.Printf("params:%T\n", params)
	queuePrefix := consts.QueuePrefix
	contain := container.CreateContainersFactory()
	// 从容器中取出一个实例
	instance := contain.Get(queuePrefix + event)
	// 将实例转换为队列处理实例
	queueInstance, ok := instance.(interf.QueueInterface)
	if !ok {
		fmt.Printf("%s is not implement QueueInterface")
	}
	queueInstance.Handle(params)
}

// Pusher 将消息投递到队列中
func Pusher(event string, data string, mode string, delay int) bool {
	// 获取队列生产者
	producer := GetQueuePusherMode(mode)
	msg := fmt.Sprintf(QueueMessageFormat, event, data)
	fmt.Println("最终发送的msg", msg)
	return producer.Send(msg, delay)
}

// GetQueuePusherMode 获取队列发布者mode
func GetQueuePusherMode(mode string) interf.QueuePusherInterface {
	var handle interf.QueuePusherInterface
	switch mode {
	case "publish_subscribe":
		handle, _ = publish_subscribe.CreateProducer()
		break
	default:
		handle, _ = publish_subscribe.CreateProducer()
		break
	}

	return handle
}
