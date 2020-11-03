package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sarama_producer/internal/handler"
	"sarama_producer/internal/kafka"
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

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8090", nil))
}
