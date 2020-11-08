package handler

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"net/http"
	"sarama_producer/internal/kafka"
	"sarama_producer/internal/logs"
	"sarama_producer/pkg/person"
)

const personContext = "Person"

type PersonHandler interface {
	Handle(responseWriter http.ResponseWriter, request *http.Request)
}

type PersonHandlerImpl struct {
	topic    string
	producer kafka.Producer
}

func NewPersonHandlerImpl(topic string, producer kafka.Producer) PersonHandler {
	return &PersonHandlerImpl{topic: topic, producer: producer}
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
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(messageKey),
		Value: sarama.ByteEncoder(body),
	}

	partition, offset, err := p.producer.Produce(message)
	if err != nil {
		logs.Logger.Error("fail to produce person",
			zap.Error(err),
			zap.String("topic", p.topic),
			zap.String("key", messageKey),
			zap.String("context", personContext),
			zap.String("lib", logs.Lib),
			zap.String("projectType", logs.ProjectType))

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	logs.Logger.Info("person produced",
		zap.String("topic", p.topic),
		zap.String("key", messageKey),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.String("context", personContext),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))
}
