#!/usr/bin/env bash

docker logs -f --tail=50 confluentconsumer | sed -e 's/^/[confluent-consumer]/ ' &
docker logs -f --tail=50 confluentproducer | sed -e 's/^/[confluent-producer]/ ' &
docker logs -f --tail=50 saramaconsumer | sed -e 's/^/[sarama-consumer]/ ' &
docker logs -f --tail=50 saramaproducer | sed -e 's/^/[sarama-producer]/ ' &
docker logs -f --tail=50 segmentioconsumer | sed -e 's/^/[segmentio-consumer]/ ' &
docker logs -f --tail=50 segmentioproducer | sed -e 's/^/[segmentio-producer]/ ' &