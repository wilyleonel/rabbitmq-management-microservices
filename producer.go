
package rabbitmqmanagement



import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Producer produce mensajes a RabbitMQ
type Producer struct {
	connectionManager *ConnectionManager
}

// NewProducer crea un nuevo productor
func NewProducer(connManager *ConnectionManager) *Producer {
	return &Producer{
		connectionManager: connManager,
	}
}

// Publish envía un mensaje
func (p *Producer) Publish(exchange, routingKey string, message interface{}, options ...PublishOption) error {
	channel, err := p.connectionManager.GetChannel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}

	// Serializar mensaje
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Configurar opciones de publicación
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
	}

	// Aplicar opciones adicionales
	for _, option := range options {
		option(&publishing)
	}

	// Publicar mensaje
	err = channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		publishing,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("📤 Message published to %s with routing key %s", exchange, routingKey)
	return nil
}

// PublishOption opción para configurar la publicación
type PublishOption func(*amqp.Publishing)

// WithMessageID establece el ID del mensaje
func WithMessageID(messageID string) PublishOption {
	return func(p *amqp.Publishing) {
		p.MessageId = messageID
	}
}

// WithHeaders establece headers personalizados
func WithHeaders(headers amqp.Table) PublishOption {
	return func(p *amqp.Publishing) {
		p.Headers = headers
	}
}

// WithExpiration establece expiración
func WithExpiration(expiration string) PublishOption {
	return func(p *amqp.Publishing) {
		p.Expiration = expiration
	}
}