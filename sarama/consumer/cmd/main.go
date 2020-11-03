package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sarama_consumer/internal/kafka"
	"sarama_consumer/pkg"
)

func main() {
	log.Println("### starting Sarama consumer ###")

	consumer, err := kafka.NewConsumer()
	if err != nil {
		panic(err)
	}

	go pkg.NewKafkaSubscriberImpl(consumer, []string{pkg.PersonTopic, pkg.TransactionTopic}).Start()

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8091", nil))
}
