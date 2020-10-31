package kafka

import (
	"github.com/Shopify/sarama"
	"sarama_producer/config"
)

type Producer interface {
	Produce(message *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type SyncProducer struct {
	sync sarama.SyncProducer
}

func NewSyncProducer() (Producer, error) {
	c := sarama.NewConfig()
	c.Version = config.KafkaVersion
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true
	asyncProducer, err := sarama.NewSyncProducer(config.KafkaBrokers(), c)
	if err != nil {
		return nil, err
	}

	return &SyncProducer{sync: asyncProducer}, nil
}

func NewSaramaSyncProducer(saramaSyncProducer sarama.SyncProducer) Producer {
	return &SyncProducer{sync: saramaSyncProducer}
}

func (p *SyncProducer) Produce(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return p.sync.SendMessage(message)
}
