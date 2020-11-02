package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"segmentio_producer/internal/handler"
	"segmentio_producer/internal/kafka"
)

const (
	PersonTopic      = "person-segmentio"
	TransactionTopic = "transaction-segmentio"
)

func main() {
	log.Println("### starting Segmentio producer ###")

	personProducer := kafka.NewSyncProducer(PersonTopic)
	http.HandleFunc("/segmentio/person", handler.NewPersonHandlerImpl(personProducer).Handle)

	transactionProducer := kafka.NewSyncProducer(TransactionTopic)
	http.HandleFunc("/segmentio/transaction", handler.NewTransactionHandlerImpl(transactionProducer).Handle)

	http.Handle("/segmentio/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8070", nil))
}
