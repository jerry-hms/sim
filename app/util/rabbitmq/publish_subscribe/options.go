package publish_subscribe

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"sim/app/global/variable"
)

type OptionsProd interface {
	apply(*producer)
}

type OptionFunc func(*producer)

func (o OptionFunc) apply(p *producer) {
	o.apply(p)
}

func SetProdMsgDelayParams(enableMsgDelayPlugin bool) OptionsProd {
	return OptionFunc(func(p *producer) {
		p.enableDelayMsgPlugin = enableMsgDelayPlugin
		p.exchangeType = "x-delayed-message"
		p.args = amqp.Table{
			"x-delayed-type": "sim",
		}
		p.exchangeName = variable.ConfigYml.GetString("rabbitMq.publishSubscribe.delayedExchangeName")
		// 延迟消息队列，交换机、消息全部设置为持久
		p.durable = true
	})
}

type OptionsConsumer interface {
	apply(*consumer)
}

type OptionsConsumerFunc func(*consumer)

func (o OptionsConsumerFunc) apply(c *consumer) {
	o.apply(c)
}

// SetConsMsgDelayParams 开发者设置消费者端初始化时的参数
func SetConsMsgDelayParams(enableDelayMsgPlugin bool) OptionsConsumer {
	return OptionsConsumerFunc(func(c *consumer) {
		c.enableDelayMsgPlugin = enableDelayMsgPlugin
		c.exchangeType = "x-delayed-message"
		c.exchangeName = variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.DelayedExchangeName")
		// 延迟消息队列，交换机、消息全部设置为持久
		c.durable = true
	})
}
