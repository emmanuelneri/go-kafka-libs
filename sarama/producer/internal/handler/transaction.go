package handler

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"net/http"
	"sarama_producer/internal/kafka"
	"sarama_producer/pkg/transaction"
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

	message := &sarama.ProducerMessage{
		Topic: t.topic,
		Key:   sarama.ByteEncoder(transactionRequested.Identifier),
		Value: sarama.ByteEncoder(body),
	}

	partition, offset, err := t.producer.Produce(message)
	if err != nil {
		log.Printf("[ERROR] fail to produce transaction. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[%s] - transaction %s produced. partition: %d - offset: %d",
		t.topic, transactionRequested.Identifier, partition, offset)
}
