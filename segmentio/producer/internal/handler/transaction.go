package handler

import (
	"context"
	"encoding/json"
	segmentioKafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"net/http"
	logs "segmentio_producer/internal"
	"segmentio_producer/internal/kafka"
	"segmentio_producer/pkg/transaction"
)

const TransactionContext = "Transaction"

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
	message := segmentioKafka.Message{
		Key:   []byte(messageKey),
		Value: body,
	}

	err = t.producer.Produce(context.Background(), message)
	if err != nil {
		logs.Logger.Error("fail to produce transaction",
			zap.Error(err),
			zap.String("topic", t.producer.Topic()),
			zap.String("key", messageKey),
			zap.String("context", TransactionContext),
			zap.String("lib", logs.Lib),
			zap.String("projectType", logs.ProjectType))

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	logs.Logger.Info("transaction produced",
		zap.String("topic", t.producer.Topic()),
		zap.String("key", messageKey),
		zap.String("context", TransactionContext),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))
}
