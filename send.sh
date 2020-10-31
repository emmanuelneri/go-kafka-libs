#!/usr/bin/env bash

URL=http://localhost:8080
TOTAL=10

START=$(date +%s)

echo "send starting... ${TOTAL} registers"
for i in $(seq 1 $TOTAL);
do
  curl --request POST \
    --url ${URL}/sarama/person \
    --header 'Content-Type: application/json' \
    --data "{\"document\": \"${i}\",\"name\": \"Customer 1\"}"

  curl --request POST \
      --url ${URL}/sarama/transaction \
      --header 'Content-Type: application/json' \
      --data "{\"identifier\": \"${i}\",\"customer\": \"Customer ${i}\",\"Value\": ${i}}"
done