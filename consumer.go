package rabbitmqmanagement



import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

// Consumer consume mensajes de RabbitMQ
type Consumer struct {
	connectionManager *ConnectionManager
	retryManager      *RetryManager
	deadLetterManager *DeadLetterManager
	handler           MessageHandler
	config            ConsumerConfig
	wg                sync.WaitGroup
	stopChan          chan struct{}
}

// NewConsumer crea un nuevo consumidor
func NewConsumer(
	connManager *ConnectionManager,
	retryManager *RetryManager,
	deadLetterManager *DeadLetterManager,
	handler MessageHandler,
	config ConsumerConfig,
) *Consumer {
	return &Consumer{
		connectionManager: connManager,
		retryManager:      retryManager,
		deadLetterManager: deadLetterManager,
		handler:           handler,
		config:            config,
		stopChan:          make(chan struct{}),
	}
}

// Start inicia el consumo de mensajes
func (c *Consumer) Start() error {
	channel, err := c.connectionManager.GetChannel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}

	// Configurar dead letter si está habilitado
	if c.deadLetterManager != nil {
		err = c.deadLetterManager.SetupDeadLetter(channel)
		if err != nil {
			return fmt.Errorf("failed to setup dead letter: %w", err)
		}
	}

	// Consumir mensajes
	msgs, err := channel.Consume(
		c.config.QueueName,
		c.config.ConsumerTag,
		c.config.AutoAck,
		c.config.Exclusive,
		c.config.NoLocal,
		c.config.NoWait,
		c.config.Args,
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	c.wg.Add(1)
	go c.processMessages(msgs)

	log.Printf("🚀 Consumer started for queue: %s", c.config.QueueName)
	return nil
}

// processMessages procesa los mensajes recibidos
func (c *Consumer) processMessages(msgs <-chan amqp.Delivery) {
	defer c.wg.Done()

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Println("📭 Channel closed, stopping consumer")
				return
			}
			c.handleMessage(msg)
		case <-c.stopChan:
			log.Println("🛑 Consumer stopped")
			return
		}
	}
}

// handleMessage maneja un mensaje individual
func (c *Consumer) handleMessage(msg amqp.Delivery) {
	err := c.handler(msg)
	
	if err != nil {
		log.Printf("❌ Error processing message: %v", err)
		c.handleError(msg, err)
	} else {
		// ✅ Éxito - confirmar mensaje
		if err := msg.Ack(false); err != nil {
			log.Printf("⚠️ Failed to ack message: %v", err)
		}
		log.Printf("✅ Message processed successfully")
	}
}

// handleError maneja errores de procesamiento
func (c *Consumer) handleError(msg amqp.Delivery, processingErr error) {
	channel, err := c.connectionManager.GetChannel()
	if err != nil {
		log.Printf("🚨 Failed to get channel for error handling: %v", err)
		return
	}

	// Intentar reintento si está configurado
	if c.retryManager != nil && c.retryManager.ShouldRetry(msg) {
		err := c.retryManager.HandleRetry(msg, channel, "retry_exchange", "retry_queue")
		if err == nil {
			// Rechazar mensaje original para que sea reencolado
			msg.Nack(false, false)
			return
		}
	}

	// Si no hay reintento o falló el reintento, enviar a DLQ
	if c.deadLetterManager != nil {
		err := c.deadLetterManager.SendToDeadLetter(msg, channel, processingErr.Error())
		if err != nil {
			log.Printf("🚨 Failed to send to DLQ: %v", err)
		}
	}

	// Rechazar mensaje original
	msg.Nack(false, false)
}

// Stop detiene el consumidor
func (c *Consumer) Stop() {
	close(c.stopChan)
	c.wg.Wait()
	log.Println("🛑 Consumer stopped gracefully")
}