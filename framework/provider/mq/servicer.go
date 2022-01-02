package mq

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type HadeMQ struct {
	addr   string         // kafka地址
	topic  string         // topic主题
	config *sarama.Config // kafka配置
}

func NewHadeMQ(params ...interface{}) (interface{}, error) {
	addr := params[0].(string)
	topic := params[1].(string)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	return HadeMQ{
		addr:   addr,
		topic:  topic,
		config: config,
	}, nil
}

func (hadeMQ HadeMQ) SendMessage(message string) error {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = hadeMQ.topic
	msg.Value = sarama.StringEncoder(message)
	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{hadeMQ.addr}, hadeMQ.config)
	if err != nil {
		return err
	}
	defer client.Close()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)

	return nil
}

func (hadeMQ HadeMQ) ConsumeMessage() error {
	consumer, err := sarama.NewConsumer([]string{hadeMQ.addr}, nil)
	if err != nil {
		return err
	}
	partitionList, err := consumer.Partitions(hadeMQ.topic) // 根据topic取到所有的分区
	if err != nil {
		return err
	}
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(hadeMQ.topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			return err
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}

	return nil
}
