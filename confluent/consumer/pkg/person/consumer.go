package person

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type MessageProcessor interface {
	Process(m *kafka.Message)
}

type MessageProcessorImpl struct {
}

func NewMessageProcessorImpl() MessageProcessor {
	return &MessageProcessorImpl{}
}

func (p MessageProcessorImpl) Process(m *kafka.Message) {
	log.Printf("[%s] - person message received. partition: %d - offset %v - key: %s - Message: %s",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset, string(m.Key), string(m.Value))

}
