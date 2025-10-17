package rabbitmqmanagement

import (
	"time"

	"github.com/streadway/amqp"
)

// Config configuración de RabbitMQ
type Config struct {
	URL           string
	PrefetchCount int
	PrefetchSize  int
}

// QueueConfig configuración de cola
type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// ExchangeConfig configuración de exchange
type ExchangeConfig struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// ConsumerConfig configuración del consumer
type ConsumerConfig struct {
	QueueName     string
	ConsumerTag   string
	AutoAck       bool
	Exclusive     bool
	NoLocal       bool
	NoWait        bool
	Args          amqp.Table
	RetryEnabled  bool
	MaxRetries    int
	RetryDelay    time.Duration
	DeadLetterEnabled bool
}

// MessageHandler maneja mensajes recibidos
type MessageHandler func(message amqp.Delivery) error

// RetryConfig configuración de reintentos
type RetryConfig struct {
	Enabled   bool
	MaxRetries int
	Delay     time.Duration
	Backoff   bool
}

// DeadLetterConfig configuración de dead letter
type DeadLetterConfig struct {
	Enabled    bool
	Exchange   string
	Queue      string
	RoutingKey string
}