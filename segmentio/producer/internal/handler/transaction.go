package handler

import (
	"context"
	"encoding/json"
	segmentioKafka "github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"segmentio_producer/internal/kafka"
	"segmentio_producer/pkg/transaction"
)

type TransactionHandler interface {
	Handle(responseWriter http.ResponseWriter, request *http.Request)
}

type TransactionHandlerImpl struct {
	producer kafka.Producer
}

func NewTransactionHandlerImpl(producer kafka.Producer) TransactionHandler {
	return &TransactionHandlerImpl{producer: producer}
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

	message := segmentioKafka.Message{
		Key:   []byte(transactionRequested.Identifier),
		Value: body,
	}

	err = t.producer.Produce(context.Background(), message)
	if err != nil {
		log.Printf("[ERROR] fail to produce transaction. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[%s] - transaction %s produced.",
		t.producer.Topic(), transactionRequested.Identifier)
}
