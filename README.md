# go-kafka-libs

Simple project to validate Golang libs for Kafka producers and consumers.

### Libs
- [Sarama](https://github.com/Shopify/sarama)
- [Confluent](https://github.com/confluentinc/confluent-kafka-go)
- [Segmentio](https://github.com/segmentio/kafka-go)

### Use Case
The use case is segregate into two applications:
- Producer: Offer two HTTP interfaces to produce information into Kafka with synchronous producer
- Consumer: Consume two different topics with Kafka consumer group

![alt tag](https://github.com/emmanuelneri/go-kafka-libs/blob/main/usecase.png?style=centerme)

### Start environment
- `./start-infra.sh `
- `./build-images.sh `
- `./start-apps.sh `

### Send data to services 
- `./send.sh `

#### Monitoring
- access `cd conf` and execute `docker-compose up`  to start monitoring tools
       
Metrics:  
1. Access Grafana
- http://localhost:3000/login
    - login: admin/admin

2. Configure prometheus data
    - Add new datasource
        -http://localhost:3000/datasources/
    - Configure prometheus URL
        - prometheus:9090        

Obs: filter by job name`{job=~".*consumer"}`

Logs:   
1. Access Kibana 
    - http://localhost:5601/app/discover
2. Create index
    - http://localhost:5601/app/discover
    - Create an index for `fluentd-*`