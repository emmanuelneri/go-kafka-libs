package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

type ConsumerHandler struct {
	ready        chan bool
	consumedChan chan sarama.ConsumerMessage
}

func newConsumerHandler() *ConsumerHandler {
	return &ConsumerHandler{
		consumedChan: make(chan sarama.ConsumerMessage),
		ready:        make(chan bool),
	}
}

func (ch *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("ConsumerHandler setup")
	close(ch.ready)
	return nil
}

func (ch *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("ConsumerHandler Cleanup")
	return nil
}

func (ch *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Printf("ConsumeClaim: %s", claim.Topic())
	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		ch.consumedChan <- *message
	}

	return nil
}
