package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
	logs "segmentio_producer/internal"
	"segmentio_producer/internal/handler"
	"segmentio_producer/internal/kafka"
)

const (
	PersonTopic      = "person-segmentio"
	TransactionTopic = "transaction-segmentio"
)

func main() {
	logs.Logger.Info("starting Segmentio producer",
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	personProducer := kafka.NewSyncProducer(PersonTopic)
	http.HandleFunc("/segmentio/person", handler.NewPersonHandlerImpl(personProducer).Handle)

	transactionProducer := kafka.NewSyncProducer(TransactionTopic)
	http.HandleFunc("/segmentio/transaction", handler.NewTransactionHandlerImpl(transactionProducer).Handle)

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8070", nil))
}
