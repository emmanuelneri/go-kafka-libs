package transaction

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"sarama_consumer/internal/logs"
)

type MessageProcessor interface {
	Process(m sarama.ConsumerMessage)
}

type MessageProcessorImpl struct {
}

func NewMessageProcessorImpl() MessageProcessor {
	return &MessageProcessorImpl{}
}

func (p MessageProcessorImpl) Process(m sarama.ConsumerMessage) {
	logs.Logger.Info("transaction message received",
		zap.String("topic", m.Topic),
		zap.Int32("partition", m.Partition),
		zap.Int64("offset", m.Offset),
		zap.String("key", string(m.Key)),
		zap.String("value", string(m.Value)),
		zap.String("context", "Transaction"),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))
}
