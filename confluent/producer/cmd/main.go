package main

import (
	"confluent_producer/internal/handler"
	"confluent_producer/internal/kafka"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

const (
	PersonTopic      = "person-confluent"
	TransactionTopic = "transaction-confluent"
)

func main() {
	log.Println("### starting Confluent producer ###")

	producer, err := kafka.NewAsyncProducer()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/confluent/person", handler.NewPersonHandlerImpl(PersonTopic, producer).Handle)
	http.HandleFunc("/confluent/transaction", handler.NewTransactionHandlerImpl(TransactionTopic, producer).Handle)

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8080", nil))
}
