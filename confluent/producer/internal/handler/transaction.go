package handler

import (
	"confluent_producer/internal/kafka"
	"confluent_producer/pkg/transaction"
	"encoding/json"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"net/http"
)

type TransactionHandler interface {
	Handle(responseWriter http.ResponseWriter, request *http.Request)
}

type TransactionHandlerImpl struct {
	topic    string
	producer kafka.Producer
}

func NewTransactionHandlerImpl(topic string, producer kafka.Producer) TransactionHandler {
	return &TransactionHandlerImpl{topic: topic, producer: producer}
}

func (t TransactionHandlerImpl) Handle(writer http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		http.Error(writer, "body required", http.StatusBadRequest)
		return
	}

	transactionRequested := &transaction.Transaction{}
	err := json.NewDecoder(request.Body).Decode(&transactionRequested)
	if err != nil {
		log.Printf("[ERROR] fail to decode transaction. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	defer request.Body.Close()
	body, err := json.Marshal(transactionRequested)
	if err != nil {
		log.Printf("[ERROR] fail to Marshal transaction. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	message := &confluentKafka.Message{
		TopicPartition: confluentKafka.TopicPartition{Topic: &t.topic, Partition: confluentKafka.PartitionAny},
		Key:            []byte(transactionRequested.Identifier),
		Value:          body,
	}

	partition, offset, err := t.producer.Produce(message)
	if err != nil {
		log.Printf("[ERROR] fail to produce transaction. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[%s] - transaction %s produced. partition: %d - offset: %s",
		t.topic, transactionRequested.Identifier, partition, offset)
}
