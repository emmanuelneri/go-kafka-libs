package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"sarama_consumer/internal/logs"
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

func (ch *ConsumerHandler) Setup(s sarama.ConsumerGroupSession) error {
	logs.Logger.Info("ConsumerHandler setup",
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	close(ch.ready)
	return nil
}

func (ch *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	logs.Logger.Warn("ConsumerHandler Cleanup",
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	return nil
}

func (ch *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logs.Logger.Info(fmt.Sprintf("ConsumeClaim: %s", claim.Topic()),
		zap.String("lib", logs.Lib),
		zap.String("projectType", logs.ProjectType))

	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		ch.consumedChan <- *message
	}

	return nil
}
