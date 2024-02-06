package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sim/app/util/queue"
	"sim/app/util/rabbitmq/publish_subscribe"
)

// BootMq 启动队列
func BootMq() {
	bootPublishSubscribe()
}

func bootPublishSubscribe() {
	consumer, err := publish_subscribe.CreateConsumer()
	if err != nil {
		fmt.Println("消费队列启动报错")
	}
	consumer.OnConnectionError(func(err *amqp.Error) {
		fmt.Println("rabbitmq重连失败，超过最大重连次数")
	})
	consumer.Received(func(receiveData string) {
		// 队列转发
		queue.Transfer(receiveData)
	})
}
