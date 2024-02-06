package publish_subscribe

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sim/app/global/variable"
	"sim/app/util/queue/interf"
	"sim/app/util/rabbitmq/execption"
)

// CreateProducer 创建一个publish生产者
func CreateProducer(options ...OptionsProd) (interf.QueuePusherInterface, error) {
	// 连接rabbitmq
	conn, err := amqp.Dial(variable.ConfigYml.GetString("rabbitMq.publishSubscribe.addr"))
	if err != nil {
		fmt.Println("rabbitmq连接报错:", err.Error())
		return nil, err
	}

	prod := &producer{
		connect:      conn,
		exchangeType: variable.ConfigYml.GetString("rabbitMq.publishSubscribe.exchangeType"),
		exchangeName: variable.ConfigYml.GetString("rabbitMq.publishSubscribe.exchangeName"),
		queueName:    variable.ConfigYml.GetString("rabbitMq.publishSubscribe.queueName"),
		durable:      variable.ConfigYml.GetBool("rabbitMq.publishSubscribe.durable"),
		args:         nil,
	}

	// 设置其他参数
	for _, v := range options {
		v.apply(prod)
	}

	return prod, nil
}

type producer struct {
	connect              *amqp.Connection
	exchangeType         string
	exchangeName         string
	queueName            string
	durable              bool
	occurError           error
	enableDelayMsgPlugin bool
	args                 amqp.Table
}

func (p *producer) Send(data string, delayMillisecond int) bool {

	ch, err := p.connect.Channel()
	p.occurError = execption.ErrorDeal(err)

	defer func() {
		_ = ch.Close()
	}()

	err = ch.ExchangeDeclare(
		p.exchangeName,
		p.exchangeType,
		p.durable,
		!p.durable,
		false,
		false,
		p.args,
	)
	p.occurError = execption.ErrorDeal(err)

	// 如果队列的声明是持久化的，那么消息也设置为持久化
	msgPersistent := amqp.Transient
	if p.durable {
		msgPersistent = amqp.Persistent
	}

	if err == nil {
		err = ch.PublishWithContext(
			context.Background(),
			p.exchangeName,
			p.queueName,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: msgPersistent,
				ContentType:  "text/plain",
				Body:         []byte(data),
				Headers: amqp.Table{
					"x-delay": delayMillisecond, // 延迟时间: 毫秒
				},
			},
		)
	}
	p.occurError = execption.ErrorDeal(err)
	if p.occurError != nil {
		return false
	} else {
		return true
	}
}

// Close 发送完毕手动关闭，这样不影响send多次发送数据
func (p *producer) Close() {
	_ = p.connect.Close()
}
