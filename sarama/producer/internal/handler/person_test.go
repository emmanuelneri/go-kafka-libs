package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"sarama_producer/internal/kafka"
	"sarama_producer/pkg/person"
	"testing"
)

func TestPersonHandler(t *testing.T) {
	body, err := json.Marshal(person.Person{
		Document: "0000",
		Name:     "Test",
	})
	assert.Nil(t, err)

	mockProducer := mocks.NewSyncProducer(t, &sarama.Config{})
	mockProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
		fmt.Println("message value: " + string(val))
		if string(body) != string(val) {
			return fmt.Errorf("expected: %s, actual: %s ", string(body), string(val))
		}
		return nil
	})

	producer := kafka.NewSaramaSyncProducer(mockProducer)
	handler := NewPersonHandlerImpl("test-person", producer)

	request := httptest.NewRequest("POST", "http://localhost:8080", bytes.NewReader(body))
	handler.Handle(httptest.NewRecorder(), request)

	assert.Nil(t, mockProducer.Close())
}
