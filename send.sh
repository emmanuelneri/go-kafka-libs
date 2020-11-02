#!/usr/bin/env bash

URL=http://localhost
TOTAL=10

START=$(date +%s)

echo "send starting... ${TOTAL} registers"
for i in $(seq 1 $TOTAL);
do
  curl --request POST \
    --url ${URL}:8080/sarama/person \
    --header 'Content-Type: application/json' \
    --data "{\"document\": \"${i}\",\"name\": \"Customer 1\"}"

  curl --request POST \
      --url ${URL:8080/sarama/transaction \
      --header 'Content-Type: application/json' \
      --data "{\"identifier\": \"${i}\",\"customer\": \"Customer ${i}\",\"Value\": ${i}}"

  curl --request POST \
    --url ${URL}:8090/confluent/person \
    --header 'Content-Type: application/json' \
    --data "{\"document\": \"${i}\",\"name\": \"Customer 1\"}"

  curl --request POST \
    --url ${URL}:8090/confluent/transaction \
    --header 'Content-Type: application/json' \
    --data "{\"identifier\": \"${i}\",\"customer\": \"Customer ${i}\",\"Value\": ${i}}"

  curl --request POST \
    --url ${URL}:8070/segmentio/person \
    --header 'Content-Type: application/json' \
    --data "{\"document\": \"${i}\",\"name\": \"Customer 1\"}"

  curl --request POST \
    --url ${URL}:8070/segmentio/transaction \
    --header 'Content-Type: application/json' \
    --data "{\"identifier\": \"${i}\",\"customer\": \"Customer ${i}\",\"Value\": ${i}}"
done