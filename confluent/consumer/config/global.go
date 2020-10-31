package config

import (
	"os"
)

var (
	localKafkaBrokers = "localhost:9092"
)

func KafkaBrokers() string {
	envValue := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if envValue == "" {
		return localKafkaBrokers
	}
	return envValue
}
