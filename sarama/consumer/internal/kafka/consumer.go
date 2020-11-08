package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"sarama_consumer/config"
	"sarama_consumer/internal/logs"
)

const consumerGroupId = "sarama-consumer"

type ConsumerGroup interface {
	Subscribe(ctx context.Context, topics []string, ready chan bool)
	ConsumedChan() map[string]chan sarama.ConsumerMessage
}

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	consumedChan  map[string]chan sarama.ConsumerMessage
}

func NewConsumer() (ConsumerGroup, error) {
	c := sarama.NewConfig()
	c.Version = config.KafkaVersion
	c.Consumer.Offsets.Initial = sarama.OffsetOldest
	c.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumer, err := sarama.NewConsumerGroup(config.KafkaBrokers(), consumerGroupId, c)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumerGroup: consumer,
		consumedChan: make(map[string]chan sarama.ConsumerMessage),
	}, nil
}

func (consumer *Consumer) ConsumedChan() map[string]chan sarama.ConsumerMessage {
	return consumer.consumedChan
}

func (consumer *Consumer) Subscribe(ctx context.Context, topics []string, ready chan bool) {
	logs.Logger.Info(fmt.Sprintf("topics subscribed %s", topics),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	go logErrors(consumer.consumerGroup.Errors())

	for _, topic := range topics {
		consumer.consumedChan[topic] = make(chan sarama.ConsumerMessage)
	}

	handler := newConsumerHandler()
	go consumer.consume(ctx, topics, handler)

	<-handler.ready
	ready <- true
	for {
		message := <-handler.consumedChan
		topicChan := consumer.consumedChan[message.Topic]
		topicChan <- message
	}
}

func (consumer *Consumer) consume(ctx context.Context, topics []string, handler *ConsumerHandler) {
	for {
		err := consumer.consumerGroup.Consume(ctx, topics, handler)
		if err != nil {
			logs.Logger.Fatal("consume group error",
				zap.Error(err),
				zap.String("lib", logs.Lib),
				zap.String("projectType", logs.ProjectType))
		}
	}
}

func logErrors(errorsChan <-chan error) {
	for err := range errorsChan {
		logs.Logger.Error("consumer error channel",
			zap.Error(err),
			zap.String("lib", logs.Lib),
			zap.String("projectType", logs.ProjectType))
	}
}
