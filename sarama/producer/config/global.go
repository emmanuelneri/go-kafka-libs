package config

import (
	"github.com/Shopify/sarama"
	"os"
	"strings"
)

var (
	localKafkaBrokers = "localhost:9092"
	KafkaVersion      = sarama.V2_5_0_0
)

func KafkaBrokers() []string {
	envValue := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if envValue == "" {
		return strings.Split(localKafkaBrokers, ",")
	}
	return strings.Split(envValue, ",")
}
