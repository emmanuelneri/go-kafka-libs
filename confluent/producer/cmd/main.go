package main

import (
	"confluent_producer/internal/handler"
	"confluent_producer/internal/kafka"
	"confluent_producer/internal/logs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
)

const (
	PersonTopic      = "person-confluent"
	TransactionTopic = "transaction-confluent"
)

func main() {
	logs.Logger.Info("starting Confluent producer",
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	producer, err := kafka.NewAsyncProducer()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/confluent/person", handler.NewPersonHandlerImpl(PersonTopic, producer).Handle)
	http.HandleFunc("/confluent/transaction", handler.NewTransactionHandlerImpl(TransactionTopic, producer).Handle)

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8080", nil))
}
