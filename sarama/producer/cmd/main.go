package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"kafka_app/internal/handler"
	"kafka_app/internal/kafka"
	"log"
	"net/http"
)

const (
	PersonTopic      = "person-sarama"
	TransactionTopic = "transaction-sarama"
)

func main() {
	log.Println("### starting Sarama producer ###")

	producer, err := kafka.NewSyncProducer()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/sarama/person", handler.NewPersonHandlerImpl(PersonTopic, producer).Handle)
	http.HandleFunc("/sarama/transaction", handler.NewTransactionHandlerImpl(TransactionTopic, producer).Handle)

	http.Handle("/sarama/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8080", nil))
}
