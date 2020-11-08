#!/usr/bin/env bash

saramaURL=http://localhost:8090/sarama
confluentURL=http://localhost:8080/confluent
segmentioURL=http://localhost:8070/segmentio

personURI=person
transactionURI=transaction

urls="${confluentURL} ${segmentioURL}"
TOTAL=100

START=$(date +%s)

echo "send starting... ${TOTAL} registers"
for i in $(seq 1 $TOTAL);do
    for url in $urls; do
      echo "${url}/${personURI}"
      curl -s --request POST \
        --url ${url}/${personURI} \
        --header 'Content-Type: application/json' \
        --data "{\"document\": \"${i}\",\"name\": \"Customer ${i}\"}" &

       echo "${url}/${transactionURI}"
      curl -s --request POST \
          --url ${url}/${transactionURI} \
          --header 'Content-Type: application/json' \
          --data "{\"identifier\": \"${i}\",\"customer\": \"Customer ${i}\",\"Value\": ${i}}" &
    done
done

wait