package pkg

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sarama_consumer/internal/kafka"
	"sarama_consumer/internal/logs"
	"sarama_consumer/pkg/person"
	"sarama_consumer/pkg/transaction"
)

const (
	PersonTopic      = "person-sarama"
	TransactionTopic = "transaction-sarama"
)

type KafkaSubscriber interface {
	Start()
}

type KafkaSubscriberImpl struct {
	topics               []string
	kafka                kafka.ConsumerGroup
	personProcessor      person.MessageProcessor
	transactionProcessor transaction.MessageProcessor
}

func NewKafkaSubscriberImpl(kafka kafka.ConsumerGroup, topics []string) KafkaSubscriber {
	return &KafkaSubscriberImpl{
		topics:               topics,
		kafka:                kafka,
		personProcessor:      person.NewMessageProcessorImpl(),
		transactionProcessor: transaction.NewMessageProcessorImpl(),
	}
}

func (k KafkaSubscriberImpl) Start() {
	consumerReady := make(chan bool)
	go k.kafka.Subscribe(context.Background(), k.topics, consumerReady)
	<-consumerReady

	consumedChan := k.kafka.ConsumedChan()

	for _, topic := range k.topics {
		t := topic

		switch t {
		case PersonTopic:
			go func() {
				for {
					personChannel := consumedChan[t]
					k.personProcessor.Process(<-personChannel)
				}
			}()

		case TransactionTopic:
			go func() {
				for {
					transactionChannel := consumedChan[t]
					k.transactionProcessor.Process(<-transactionChannel)
				}
			}()
		default:
			logs.Logger.Panic(fmt.Sprintf("topic not mapped %s", t),
				zap.String("lib", logs.Lib),
				zap.String("projectType", logs.ProjectType))
		}
	}
}
