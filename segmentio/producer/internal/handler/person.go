package handler

import (
	"context"
	"encoding/json"
	segmentioKafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"net/http"
	logs "segmentio_producer/internal"
	"segmentio_producer/internal/kafka"
	"segmentio_producer/pkg/person"
)

const personContext = "Person"

type PersonHandler interface {
	Handle(responseWriter http.ResponseWriter, request *http.Request)
}

type PersonHandlerImpl struct {
	producer kafka.Producer
}

func NewPersonHandlerImpl(producer kafka.Producer) PersonHandler {
	return &PersonHandlerImpl{producer: producer}
}

func (p PersonHandlerImpl) Handle(writer http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		http.Error(writer, "body required", http.StatusBadRequest)
		return
	}

	personRequested := &person.Person{}
	err := json.NewDecoder(request.Body).Decode(&personRequested)
	if err != nil {
		logs.Logger.Error("fail to decode person",
			zap.Error(err),
			zap.String("url", request.RequestURI),
			zap.String("method", request.Method),
			zap.String("context", personContext),
			zap.String("lib", logs.Lib),
			zap.String("projectType", logs.ProjectType))

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	defer request.Body.Close()
	body, err := json.Marshal(personRequested)
	if err != nil {
		logs.Logger.Error("fail to Marshal person",
			zap.Error(err),
			zap.String("url", request.RequestURI),
			zap.String("method", request.Method),
			zap.String("context", personContext),
			zap.String("lib", logs.Lib),
			zap.String("projectType", logs.ProjectType))

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	messageKey := personRequested.Document
	message := segmentioKafka.Message{
		Key:   []byte(messageKey),
		Value: body,
	}

	err = p.producer.Produce(context.Background(), message)
	if err != nil {
		logs.Logger.Error("fail to produce person",
			zap.Error(err),
			zap.String("topic", p.producer.Topic()),
			zap.String("key", messageKey),
			zap.String("context", personContext),
			zap.String("lib", logs.Lib),
			zap.String("projectType", logs.ProjectType))

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	logs.Logger.Info("person produced",
		zap.String("topic", p.producer.Topic()),
		zap.String("key", messageKey),
		zap.String("context", personContext),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))
}
