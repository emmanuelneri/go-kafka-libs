package person

import (
	"github.com/segmentio/kafka-go"
	"log"
)

type MessageProcessor interface {
	Process(m kafka.Message)
}

type MessageProcessorImpl struct {
}

func NewMessageProcessorImpl() MessageProcessor {
	return &MessageProcessorImpl{}
}

func (p MessageProcessorImpl) Process(m kafka.Message) {
	log.Printf("[%s] - person message received. partition: %d - offset %v - key: %s - Message: %s",
		m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

}
