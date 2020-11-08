package pkg

import (
	"confluent_consumer/internal/kafka"
	"confluent_consumer/internal/logs"
	"confluent_consumer/pkg/person"
	"confluent_consumer/pkg/transaction"
	"fmt"
	"go.uber.org/zap"
)

const (
	PersonTopic      = "person-confluent"
	TransactionTopic = "transaction-confluent"
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
	k.kafka.Subscribe(k.topics)
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
