package kafka

import (
	"confluent_consumer/config"
	"confluent_consumer/internal/logs"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const consumerGroupId = "confluent-consumer"

type ConsumerGroup interface {
	Subscribe(topics []string)
	ConsumedChan() map[string]chan *kafka.Message
}

type Consumer struct {
	kafkaConsumer *kafka.Consumer
	consumedChan  map[string]chan *kafka.Message
}

func NewConsumer() (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBrokers(),
		"group.id":          consumerGroupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	return &Consumer{kafkaConsumer: consumer,
		consumedChan: make(map[string]chan *kafka.Message),
	}, nil
}

func (c *Consumer) Subscribe(topics []string) {
	logs.Logger.Info(fmt.Sprintf("topics subscribed %s", topics),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	for _, topic := range topics {
		c.consumedChan[topic] = make(chan *kafka.Message)
	}

	err := c.kafkaConsumer.SubscribeTopics(topics, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			eventConsumed := c.kafkaConsumer.Poll(0)
			switch event := eventConsumed.(type) {
			case *kafka.Message:
				topicChan := c.consumedChan[*event.TopicPartition.Topic]
				topicChan <- event
			case kafka.PartitionEOF:
				logs.Logger.Error("Consume partitionEOF event",
					zap.Error(err),
					zap.Any("event", event),
					zap.String("lib", logs.Lib),
					zap.String("projectType", logs.ProjectType))

			case kafka.Error:
				logs.Logger.Fatal("Consume kafka error event",
					zap.Error(err),
					zap.Any("event", event),
					zap.String("lib", logs.Lib),
					zap.String("projectType", logs.ProjectType))
			default:
			}
		}
	}()
}

func (c *Consumer) ConsumedChan() map[string]chan *kafka.Message {
	return c.consumedChan
}
