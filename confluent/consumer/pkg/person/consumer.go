package person

import (
	"confluent_consumer/internal/logs"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
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
	logs.Logger.Info("person message received",
		zap.Stringp("topic", m.TopicPartition.Topic),
		zap.Int32("partition", m.TopicPartition.Partition),
		zap.Any("offset", m.TopicPartition.Offset),
		zap.String("key", string(m.Key)),
		zap.String("value", string(m.Value)),
		zap.String("context", "Person"),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

}
