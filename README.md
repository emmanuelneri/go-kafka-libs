# go-kafka-libs

Simple project to validate Golang libs for Kafka producers and consumers.

### Libs
- [Sarama](https://github.com/Shopify/sarama)
- [Confluent](https://github.com/confluentinc/confluent-kafka-go)
- [Segmentio](https://github.com/segmentio/kafka-go)

### Start environment
- `./start-infra.sh `
- `./build-images.sh `
- `./start-apps.sh `

### Send data to services 
- `./send.sh `

#### View services metrics
1. Access Grafana
- http://localhost:3000/login
    - default login: admin/admin

2. Configure prometheus data
    - Add new datasource
        -http://localhost:3000/datasources/
    - Configure prometheus URL
        - prometheus:9090
        
Obs: filter by job name`{job=~".*consumer"}`