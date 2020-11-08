package transaction

import (
	"confluent_consumer/internal/logs"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
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
	logs.Logger.Info("transaction message received",
		zap.Stringp("topic", m.TopicPartition.Topic),
		zap.Int32("partition", m.TopicPartition.Partition),
		zap.Any("offset", m.TopicPartition.Offset),
		zap.String("key", string(m.Key)),
		zap.String("value", string(m.Value)),
		zap.String("context", "Transaction"),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))
}
