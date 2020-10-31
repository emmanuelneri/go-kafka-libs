package handler

import (
	"confluent_producer/internal/kafka"
	"confluent_producer/pkg/person"
	"encoding/json"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"net/http"
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

	message := &confluentKafka.Message{
		TopicPartition: confluentKafka.TopicPartition{Topic: &p.topic, Partition: confluentKafka.PartitionAny},
		Key:            []byte(personRequested.Document),
		Value:          body,
	}

	partition, offset, err := p.producer.Produce(message)
	if err != nil {
		log.Printf("[ERROR] fail to produce person. %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[%s] - person %s produced. partition: %d - offset: %s",
		p.topic, personRequested.Document, partition, offset)
}
