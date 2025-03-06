package main

import (
	"lgm8-notification-service/config"
	"lgm8-notification-service/internal/email"
	"lgm8-notification-service/internal/handlers"
	"lgm8-notification-service/internal/rabbitmq"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: [%s]", err)
	}

	// Initialize Email Sender
	emailSender := email.NewEmailSender(
		cfg.SMTP.Enabled,
		cfg.SMTP.URL,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.From,
	)

	// Initialize RabbitMQ consumer manager
	consumerManager, err := rabbitmq.NewConsumerManager(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: [%s]", err)
	}
	defer consumerManager.Close()

	// Register queues and their handlers
	err = consumerManager.RegisterQueue("user-verification-email", handlers.NewUserVerificationEmailHandler(emailSender))
	if err != nil {
		log.Fatalf("Failed to register queue: [%s]", err)
	}

	// Start listening for messages
	if err := consumerManager.StartListening(); err != nil {
		log.Fatalf("Failed to start listening: [%s]", err)
	}

	// Block the main thread
	select {}
}
