package transaction

import (
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"segmentio_consumer/internal/logs"
)

type MessageProcessor interface {
	Process(kafka.Message)
}

type MessageProcessorImpl struct {
}

func NewMessageProcessorImpl() MessageProcessor {
	return &MessageProcessorImpl{}
}

func (p MessageProcessorImpl) Process(m kafka.Message) {
	logs.Logger.Info("transaction message received",
		zap.String("topic", m.Topic),
		zap.Int("partition", m.Partition),
		zap.Int64("offset", m.Offset),
		zap.String("key", string(m.Key)),
		zap.String("value", string(m.Value)),
		zap.String("context", "Transaction"),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))
}
