package person

import (
	"github.com/Shopify/sarama"
	"log"
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
	log.Printf("[%s] - person message received. partition: %d - offset %d - key: %s - Message: %s",
		m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

}
