package handler

import (
	"confluent_producer/internal/kafka"
	"confluent_producer/internal/logs"
	"confluent_producer/pkg/transaction"
	"encoding/json"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"net/http"
)

const TransactionContext = "Transaction"

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
		logs.Logger.Error("fail to decode transaction",
			zap.Error(err),
			zap.String("url", request.RequestURI),
			zap.String("method", request.Method),
			zap.String("context", TransactionContext),
			zap.String("lib", logs.Lib),
			zap.String("projectType", logs.ProjectType))

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	defer request.Body.Close()
	body, err := json.Marshal(transactionRequested)
	if err != nil {
		logs.Logger.Error("fail to Marshal transaction",
			zap.Error(err),
			zap.String("url", request.RequestURI),
			zap.String("method", request.Method),
			zap.String("context", TransactionContext),
			zap.String("lib", logs.Lib),
			zap.String("projectType", logs.ProjectType))

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	messageKey := transactionRequested.Identifier
	message := &confluentKafka.Message{
		TopicPartition: confluentKafka.TopicPartition{Topic: &t.topic, Partition: confluentKafka.PartitionAny},
		Key:            []byte(messageKey),
		Value:          body,
	}

	partition, offset, err := t.producer.Produce(message)
	if err != nil {
		logs.Logger.Error("fail to produce transaction",
			zap.Error(err),
			zap.String("topic", t.topic),
			zap.String("key", messageKey),
			zap.String("context", TransactionContext),
			zap.String("lib", logs.Lib),
			zap.String("projectType", logs.ProjectType))

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	logs.Logger.Info("transaction produced",
		zap.String("topic", t.topic),
		zap.String("key", messageKey),
		zap.Int32("partition", partition),
		zap.String("offset", offset),
		zap.String("context", TransactionContext),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))
}
