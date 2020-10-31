package kafka

import (
	"confluent_consumer/config"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
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
	fmt.Printf("Subscribe %s", topics)

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
				fmt.Printf("PartitionEOF %v\n", event)
			case kafka.Error:
				log.Panicf("Message error Error: %v\n", event)
			default:
			}
		}
	}()
}

func (c *Consumer) ConsumedChan() map[string]chan *kafka.Message {
	return c.consumedChan
}
