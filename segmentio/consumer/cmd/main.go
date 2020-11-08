package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
	"segmentio_consumer/internal/logs"
	"segmentio_consumer/pkg"
)

func main() {
	logs.Logger.Info("starting Segmentio consumer",
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	go pkg.NewKafkaSubscriberImpl([]string{pkg.PersonTopic, pkg.TransactionTopic}).Start()

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8071", nil))
}
