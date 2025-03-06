package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer manages the connection to RabbitMQ and message consumption
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewConsumer initializes a RabbitMQ connection and channel
func NewConsumer(rabbitMQURL string) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: [%w]", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: [%w]", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
	}, nil
}

// DeclareQueue ensures the queue exists before consuming messages
func (c *Consumer) DeclareQueue(queueName string) error {
	_, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue [%s]: [%w]", queueName, err)
	}
	return nil
}

// ConsumeMessages starts listening for messages on the specified queue
func (c *Consumer) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: [%w]", err)
	}

	log.Printf("Listening for messages on queue: [%s]", queueName)
	return msgs, nil
}

// Close closes the RabbitMQ connection and channel
func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
