package main

import (
	"confluent_consumer/internal/kafka"
	"confluent_consumer/pkg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	log.Println("### starting Confluent consumer ###")

	consumer, err := kafka.NewConsumer()
	if err != nil {
		panic(err)
	}

	go pkg.NewKafkaSubscriberImpl(consumer, []string{pkg.PersonTopic, pkg.TransactionTopic}).Start()

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8081", nil))
}
