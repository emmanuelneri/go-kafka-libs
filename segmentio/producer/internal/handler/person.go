package handler

import (
	"context"
	"encoding/json"
	segmentioKafka "github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"segmentio_producer/internal/kafka"
	"segmentio_producer/pkg/person"
)

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

	message := segmentioKafka.Message{
		Key:   []byte(personRequested.Document),
		Value: body,
	}

	err = p.producer.Produce(context.Background(), message)
	if err != nil {
		log.Printf("[ERROR] fail to produce person. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[%s] - person %s produced.",
		p.producer.Topic(), personRequested.Document)
}
