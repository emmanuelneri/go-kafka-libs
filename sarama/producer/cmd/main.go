package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
	"sarama_producer/internal/handler"
	"sarama_producer/internal/kafka"
	"sarama_producer/internal/logs"
)

const (
	PersonTopic      = "person-sarama"
	TransactionTopic = "transaction-sarama"
)

func main() {
	logs.Logger.Info("starting Sarama producer",
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	producer, err := kafka.NewSyncProducer()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/sarama/person", handler.NewPersonHandlerImpl(PersonTopic, producer).Handle)
	http.HandleFunc("/sarama/transaction", handler.NewTransactionHandlerImpl(TransactionTopic, producer).Handle)

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8090", nil))
}
