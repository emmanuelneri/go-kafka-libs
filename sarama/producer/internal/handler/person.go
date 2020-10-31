package handler

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"net/http"
	"sarama_producer/internal/kafka"
	"sarama_producer/pkg/person"
)

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
		log.Printf("[ERROR] fail to decode person. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	defer request.Body.Close()
	body, err := json.Marshal(personRequested)
	if err != nil {
		log.Printf("[ERROR] fail to Marshal person. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(personRequested.Document),
		Value: sarama.ByteEncoder(body),
	}

	partition, offset, err := p.producer.Produce(message)
	if err != nil {
		log.Printf("[ERROR] fail to produce person. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[%s] - person %s produced. partition: %d - offset: %d",
		p.topic, personRequested.Document, partition, offset)
}
