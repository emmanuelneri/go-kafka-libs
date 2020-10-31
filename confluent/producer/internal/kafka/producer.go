package kafka

import (
	"confluent_producer/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
)

type Producer interface {
	Produce(message *kafka.Message) (partition int32, offset string, err error)
}

type SyncProducer struct {
	producer *kafka.Producer
}

func NewAsyncProducer() (Producer, error) {
	bootstrapServers := config.KafkaBrokers()
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return nil, err
	}

	return &SyncProducer{producer: producer}, nil
}

func (p *SyncProducer) Produce(message *kafka.Message) (partition int32, offset string, err error) {
	deliveryChan := make(chan kafka.Event)
	err = p.producer.Produce(message, deliveryChan)
	if err != nil {
		close(deliveryChan)
		return 0, "", errors.Wrap(err, "producer error")
	}

	e := <-deliveryChan
	close(deliveryChan)
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return 0, "", errors.Wrap(m.TopicPartition.Error, "delivery error")
	}

	return m.TopicPartition.Partition, m.TopicPartition.Offset.String(), nil
}
