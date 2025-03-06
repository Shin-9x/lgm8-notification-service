package rabbitmq

import (
	"fmt"
	"log"
)

// MessageHandler defines an interface that each handler must implement
type MessageHandler interface {
	HandleMessage(body []byte) error
}

// ConsumerManager manages multiple queues and their handlers
type ConsumerManager struct {
	consumer *Consumer
	handlers map[string]MessageHandler
}

// NewConsumerManager initializes the ConsumerManager
func NewConsumerManager(rabbitMQURL string) (*ConsumerManager, error) {
	consumer, err := NewConsumer(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	return &ConsumerManager{
		consumer: consumer,
		handlers: make(map[string]MessageHandler),
	}, nil
}

// RegisterQueue associates a queue with a specific handler
func (m *ConsumerManager) RegisterQueue(queueName string, handler MessageHandler) error {
	if err := m.consumer.DeclareQueue(queueName); err != nil {
		return err
	}
	m.handlers[queueName] = handler
	return nil
}

// StartListening listens to all registered queues and routes messages to the appropriate handlers
func (m *ConsumerManager) StartListening() error {
	for queue, handler := range m.handlers {
		msgs, err := m.consumer.ConsumeMessages(queue)
		if err != nil {
			return fmt.Errorf("error starting consumer for queue [%s]: [%w]", queue, err)
		}

		go func(queue string, handler MessageHandler) {
			for msg := range msgs {
				log.Printf("Received message from queue [%s]", queue)
				if err := handler.HandleMessage(msg.Body); err != nil {
					log.Printf("Error handling message from queue [%s]: [%s]", queue, err)
				}
			}
		}(queue, handler)
	}
	return nil
}

// Close shuts down the consumer
func (m *ConsumerManager) Close() {
	m.consumer.Close()
}
