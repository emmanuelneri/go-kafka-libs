package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"segmentio_producer/config"
)

type Producer interface {
	Produce(ctx context.Context, message kafka.Message) error
	Topic() string
}

type SyncProducer struct {
	writer *kafka.Writer
}

func NewSyncProducer(topic string) Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.KafkaBrokers()),
		Topic:        topic,
		RequiredAcks: kafka.RequireOne,
	}

	return &SyncProducer{writer: writer}
}

func (p *SyncProducer) Produce(ctx context.Context, message kafka.Message) error {
	return p.writer.WriteMessages(ctx, message)
}

func (p SyncProducer) Topic() string {
	return p.writer.Topic
}
