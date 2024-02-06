package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sim/app/global/variable"
	"sim/app/util/queue"
	"sim/app/util/rabbitmq/publish_subscribe"
)

// InitRabbitMq 启动队列
func InitRabbitMq() {
	bootPublishSubscribe()
}

// 启动生产者消费者模式
func bootPublishSubscribe() {
	consumer, err := publish_subscribe.CreateConsumer()
	if err != nil {
		fmt.Println("消费队列启动报错")
	}
	consumer.OnConnectionError(func(err *amqp.Error) {
		variable.ZapLog.Error("")
	})
	consumer.Received(func(receiveData string) {
		// 队列转发
		queue.Transfer(receiveData)
	})
}
