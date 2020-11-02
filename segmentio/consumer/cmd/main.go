package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"segmentio_consumer/pkg"
)

func main() {
	log.Println("### starting Segmentio consumer ###")

	go pkg.NewKafkaSubscriberImpl([]string{pkg.PersonTopic, pkg.TransactionTopic}).Start()

	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(http.ListenAndServe(":8071", nil))
}
