package rabbitmqmanagement


import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// RetryManager gestiona reintentos de mensajes
type RetryManager struct {
	config RetryConfig
}

// NewRetryManager crea un nuevo gestor de reintentos
func NewRetryManager(config RetryConfig) *RetryManager {
	return &RetryManager{
		config: config,
	}
}

// ShouldRetry determina si se debe reintentar
func (rm *RetryManager) ShouldRetry(delivery amqp.Delivery) bool {
	if !rm.config.Enabled {
		return false
	}

	retryCount := rm.getRetryCount(delivery)
	return retryCount < rm.config.MaxRetries
}

// GetRetryDelay calcula el delay para el reintento
func (rm *RetryManager) GetRetryDelay(delivery amqp.Delivery) time.Duration {
	if !rm.config.Backoff {
		return rm.config.Delay
	}

	retryCount := rm.getRetryCount(delivery)
	// Backoff exponencial
	delay := rm.config.Delay * time.Duration(1<<uint(retryCount))
	
	// Limitar el delay máximo a 5 minutos
	if delay > 5*time.Minute {
		delay = 5 * time.Minute
	}
	
	return delay
}

// PrepareForRetry prepara el mensaje para reintento
func (rm *RetryManager) PrepareForRetry(delivery *amqp.Delivery) {
	retryCount := rm.getRetryCount(*delivery) + 1
	
	if delivery.Headers == nil {
		delivery.Headers = amqp.Table{}
	}
	
	delivery.Headers["x-retry-count"] = retryCount
	delivery.Headers["x-retry-delay"] = rm.GetRetryDelay(*delivery).String()
	delivery.Headers["x-last-retry"] = time.Now().Format(time.RFC3339)
}

func (rm *RetryManager) getRetryCount(delivery amqp.Delivery) int {
	if delivery.Headers == nil {
		return 0
	}
	
	if retryCount, ok := delivery.Headers["x-retry-count"].(int32); ok {
		return int(retryCount)
	}
	
	if retryCount, ok := delivery.Headers["x-retry-count"].(int); ok {
		return retryCount
	}
	
	return 0
}

// HandleRetry maneja el reintento del mensaje
func (rm *RetryManager) HandleRetry(delivery amqp.Delivery, channel *amqp.Channel, retryExchange, retryQueue string) error {
	if !rm.ShouldRetry(delivery) {
		return fmt.Errorf("max retries exceeded")
	}

	rm.PrepareForRetry(&delivery)
	
	delay := rm.GetRetryDelay(delivery)
	
	// Publicar a cola de retraso
	err := channel.Publish(
		retryExchange,
		retryQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  delivery.ContentType,
			Body:         delivery.Body,
			Headers:      delivery.Headers,
			DeliveryMode: delivery.DeliveryMode,
			MessageId:    delivery.MessageId,
			Timestamp:    time.Now(),
			Expiration:   fmt.Sprintf("%d", delay.Milliseconds()),
		},
	)
	
	if err != nil {
		return fmt.Errorf("failed to publish to retry queue: %w", err)
	}
	
	log.Printf("🔄 Message scheduled for retry in %v (attempt %d)", delay, rm.getRetryCount(delivery)+1)
	return nil
}