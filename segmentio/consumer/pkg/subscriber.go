package pkg

import (
	"context"
	"fmt"
	"log"
	"segmentio_consumer/internal/kafka"
	"segmentio_consumer/pkg/person"
	"segmentio_consumer/pkg/transaction"
)

const (
	PersonTopic      = "person-segmentio"
	TransactionTopic = "transaction-segmentio"
)

type KafkaSubscriber interface {
	Start()
}

type KafkaSubscriberImpl struct {
	topics               []string
	personProcessor      person.MessageProcessor
	transactionProcessor transaction.MessageProcessor
}

func NewKafkaSubscriberImpl(topics []string) KafkaSubscriber {
	return &KafkaSubscriberImpl{
		topics:               topics,
		personProcessor:      person.NewMessageProcessorImpl(),
		transactionProcessor: transaction.NewMessageProcessorImpl(),
	}
}

func (k KafkaSubscriberImpl) Start() {

	for _, topic := range k.topics {
		t := topic

		switch t {
		case PersonTopic:
			go func() {
				consumer := kafka.NewConsumer(t)
				for {
					message, err := consumer.FetchMessage(context.Background())
					if err != nil {
						log.Printf("[ERROR] - fail to consume topic: %s. %s", t, err)
						continue
					}

					k.personProcessor.Process(message)
				}
			}()

		case TransactionTopic:
			go func() {
				consumer := kafka.NewConsumer(t)
				for {
					message, err := consumer.FetchMessage(context.Background())
					if err != nil {
						log.Printf("[ERROR] - fail to consume topic: %s. %s", t, err)
						continue
					}

					k.transactionProcessor.Process(message)
				}
			}()
		default:
			panic(fmt.Sprintf("topic not mapped %s", t))
		}
	}
}
