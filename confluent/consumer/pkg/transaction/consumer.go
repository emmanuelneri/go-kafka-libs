package transaction

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type MessageProcessor interface {
	Process(*kafka.Message)
}

type MessageProcessorImpl struct {
}

func NewMessageProcessorImpl() MessageProcessor {
	return &MessageProcessorImpl{}
}

func (p MessageProcessorImpl) Process(m *kafka.Message) {
	log.Printf("[%s] - transaction message received. partition: %d - offset %d - key: %s - Message: %s",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset, string(m.Key), string(m.Value))
}
