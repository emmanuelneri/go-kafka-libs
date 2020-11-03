#!/usr/bin/env bash

(cd confluent/consumer; docker build -t go-kafka-libs/confluent-consumer .)
(cd confluent/producer; docker build -t go-kafka-libs/confluent-producer .)

(cd sarama/consumer; docker build -t go-kafka-libs/sarama-consumer .)
(cd sarama/producer; docker build -t go-kafka-libs/sarama-producer .)

(cd segmentio/consumer; docker build -t go-kafka-libs/segmentio-consumer .)
(cd segmentio/producer; docker build -t go-kafka-libs/segmentio-producer .)