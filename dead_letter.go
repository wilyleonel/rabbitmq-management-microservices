package rabbitmqmanagement


import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// DeadLetterManager gestiona dead letter queue
type DeadLetterManager struct {
	config DeadLetterConfig
}

// NewDeadLetterManager crea un nuevo gestor de DLQ
func NewDeadLetterManager(config DeadLetterConfig) *DeadLetterManager {
	return &DeadLetterManager{
		config: config,
	}
}

// SetupDeadLetter configura la dead letter queue
func (dlm *DeadLetterManager) SetupDeadLetter(channel *amqp.Channel) error {
	if !dlm.config.Enabled {
		return nil
	}

	// Declarar exchange de dead letter
	err := channel.ExchangeDeclare(
		dlm.config.Exchange,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLX: %w", err)
	}

	// Declarar cola de dead letter
	_, err = channel.QueueDeclare(
		dlm.config.Queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Binding entre exchange y cola
	err = channel.QueueBind(
		dlm.config.Queue,
		dlm.config.RoutingKey,
		dlm.config.Exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind DLQ: %w", err)
	}

	log.Println("📨 Dead letter queue configured")
	return nil
}

// SendToDeadLetter envía mensaje a DLQ
func (dlm *DeadLetterManager) SendToDeadLetter(delivery amqp.Delivery, channel *amqp.Channel, reason string) error {
	if !dlm.config.Enabled {
		return nil
	}

	// Agregar headers con información del error
	if delivery.Headers == nil {
		delivery.Headers = amqp.Table{}
	}
	
	delivery.Headers["x-dead-letter-reason"] = reason
	delivery.Headers["x-dead-letter-timestamp"] = time.Now().Format(time.RFC3339)

	err := channel.Publish(
		dlm.config.Exchange,
		dlm.config.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  delivery.ContentType,
			Body:         delivery.Body,
			Headers:      delivery.Headers,
			DeliveryMode: delivery.DeliveryMode,
			MessageId:    delivery.MessageId,
			Timestamp:    time.Now(),
		},
	)
	
	if err != nil {
		return fmt.Errorf("failed to send to DLQ: %w", err)
	}
	
	log.Printf("💀 Message sent to DLQ: %s", reason)
	return nil
}