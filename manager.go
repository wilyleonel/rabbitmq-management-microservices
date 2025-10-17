package rabbitmqmanagement


import (
	"fmt"
	"log"
)

// RabbitMQManager gestiona toda la interacción con RabbitMQ
type RabbitMQManager struct {
	connectionManager *ConnectionManager
	producer          *Producer
}

// NewRabbitMQManager crea un nuevo gestor de RabbitMQ
func NewRabbitMQManager(config Config) (*RabbitMQManager, error) {
	connManager := NewConnectionManager(config)
	
	err := connManager.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	producer := NewProducer(connManager)

	return &RabbitMQManager{
		connectionManager: connManager,
		producer:          producer,
	}, nil
}

// CreateQueue crea una cola
func (rm *RabbitMQManager) CreateQueue(config QueueConfig) error {
	channel, err := rm.connectionManager.GetChannel()
	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Printf("📝 Queue created: %s", config.Name)
	return nil
}

// CreateExchange crea un exchange
func (rm *RabbitMQManager) CreateExchange(config ExchangeConfig) error {
	channel, err := rm.connectionManager.GetChannel()
	if err != nil {
		return err
	}

	err = channel.ExchangeDeclare(
		config.Name,
		config.Type,
		config.Durable,
		config.AutoDelete,
		config.Internal,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	log.Printf("📝 Exchange created: %s", config.Name)
	return nil
}

// BindQueue bindea una cola a un exchange
func (rm *RabbitMQManager) BindQueue(queueName, exchangeName, routingKey string) error {
	channel, err := rm.connectionManager.GetChannel()
	if err != nil {
		return err
	}

	err = channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	log.Printf("🔗 Queue %s bound to exchange %s with key %s", queueName, exchangeName, routingKey)
	return nil
}

// GetProducer retorna el productor
func (rm *RabbitMQManager) GetProducer() *Producer {
	return rm.producer
}

// CreateConsumer crea un consumidor
func (rm *RabbitMQManager) CreateConsumer(
	handler MessageHandler,
	consumerConfig ConsumerConfig,
	retryConfig RetryConfig,
	deadLetterConfig DeadLetterConfig,
) *Consumer {

	var retryManager *RetryManager
	if retryConfig.Enabled {
		retryManager = NewRetryManager(retryConfig)
	}

	var deadLetterManager *DeadLetterManager
	if deadLetterConfig.Enabled {
		deadLetterManager = NewDeadLetterManager(deadLetterConfig)
	}

	return NewConsumer(
		rm.connectionManager,
		retryManager,
		deadLetterManager,
		handler,
		consumerConfig,
	)
}

// Close cierra todas las conexiones
func (rm *RabbitMQManager) Close() {
	rm.connectionManager.Close()
}