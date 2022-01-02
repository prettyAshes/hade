package contact

const MQKey = "hade:mq"

type MQ interface {
	// SendMessage 发送消息
	SendMessage(msg string) error
	// ConsumeMessage 消费消息
	ConsumeMessage() error
}
