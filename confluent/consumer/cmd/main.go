package main

import (
	"confluent_consumer/internal/kafka"
	"confluent_consumer/internal/logs"
	"confluent_consumer/pkg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	logs.Logger.Info("starting Confluent consumer",
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	consumer, err := kafka.NewConsumer()
	if err != nil {
		panic(err)
	}

	go pkg.NewKafkaSubscriberImpl(consumer, []string{pkg.PersonTopic, pkg.TransactionTopic}).Start()

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8081", nil))
}
