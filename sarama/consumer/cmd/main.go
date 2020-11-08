package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
	"sarama_consumer/internal/kafka"
	"sarama_consumer/internal/logs"
	"sarama_consumer/pkg"
)

func main() {
	logs.Logger.Info("starting Sarama consumer",
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	consumer, err := kafka.NewConsumer()
	if err != nil {
		panic(err)
	}

	go pkg.NewKafkaSubscriberImpl(consumer, []string{pkg.PersonTopic, pkg.TransactionTopic}).Start()

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8091", nil))
}
