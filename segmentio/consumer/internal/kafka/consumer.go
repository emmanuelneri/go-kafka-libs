package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"segmentio_consumer/config"
	"strings"
	"time"
)

const consumerGroupId = "segmentio-consumer"

type ConsumerGroup interface {
	Subscribe(topics []string)
	FetchMessage(ctx context.Context) (kafka.Message, error)
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(topic string) *Consumer {
	dialer := &kafka.Dialer{
		Timeout: 10 * time.Second,
	}

	conf := kafka.ReaderConfig{
		Brokers:     strings.Split(config.KafkaBrokers(), ","),
		GroupID:     consumerGroupId,
		Topic:       topic,
		Dialer:      dialer,
		StartOffset: kafka.FirstOffset,
		MinBytes:    10e3,
		MaxBytes:    10e6,
	}

	reader := kafka.NewReader(conf)
	return &Consumer{reader: reader}
}

func (c *Consumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}
