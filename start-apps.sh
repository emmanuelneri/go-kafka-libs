#!/usr/bin/env bash

docker run -itd \
    --name confluentconsumer \
    --network=go-kafka-libs \
    --link kafka \
    --log-driver=fluentd \
    --log-opt fluentd-address=localhost:24224 \
    -p 8081:8081 \
    go-kafka-libs/confluent-consumer

docker run -itd \
    --name confluentproducer \
    --network=go-kafka-libs \
    --link kafka \
    --log-driver=fluentd \
    --log-opt fluentd-address=localhost:24224 \
    -p 8080:8080 \
    go-kafka-libs/confluent-producer

docker run -itd \
    --name saramaconsumer \
    --network=go-kafka-libs \
    --link kafka \
    --log-driver=fluentd \
    --log-opt fluentd-address=localhost:24224 \
    -p 8091:8091 \
    go-kafka-libs/sarama-consumer

docker run -itd \
    --name saramaproducer \
    --network=go-kafka-libs \
    --link kafka \
    --log-driver=fluentd \
    --log-opt fluentd-address=localhost:24224 \
    -p 8090:8090 \
    go-kafka-libs/sarama-producer

docker run -itd \
    --name segmentioconsumer \
    --network=go-kafka-libs \
    --link kafka \
    --log-driver=fluentd \
    --log-opt fluentd-address=localhost:24224 \
    -p 8071:8071 \
    go-kafka-libs/segmentio-consumer

docker run -itd \
    --name segmentioproducer \
    --network=go-kafka-libs \
    --link kafka \
    --log-driver=fluentd \
    --log-opt fluentd-address=localhost:24224 \
    -p 8070:8070 \
    go-kafka-libs/segmentio-producer

docker ps