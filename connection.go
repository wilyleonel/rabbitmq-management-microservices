package rabbitmqmanagement


import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// ConnectionManager gestiona conexiones a RabbitMQ
type ConnectionManager struct {
	config     Config
	connection *amqp.Connection
	channel    *amqp.Channel
	mutex      sync.RWMutex
	connected  bool
}

// NewConnectionManager crea un nuevo gestor de conexiones
func NewConnectionManager(config Config) *ConnectionManager {
	return &ConnectionManager{
		config: config,
	}
}

// Connect conecta a RabbitMQ
func (cm *ConnectionManager) Connect() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.connected {
		return nil
	}

	var err error
	cm.connection, err = amqp.Dial(cm.config.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	cm.channel, err = cm.connection.Channel()
	if err != nil {
		cm.connection.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Configurar QoS
	err = cm.channel.Qos(
		cm.config.PrefetchCount,
		cm.config.PrefetchSize,
		false,
	)
	if err != nil {
		cm.channel.Close()
		cm.connection.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	cm.connected = true
	log.Println("✅ Connected to RabbitMQ")
	return nil
}



// GetChannel retorna el canal actual
func (cm *ConnectionManager) GetChannel() (*amqp.Channel, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if !cm.connected || cm.channel == nil {
		return nil, fmt.Errorf("not connected to RabbitMQ")
	}

	return cm.channel, nil
}

// IsConnected verifica si está conectado
func (cm *ConnectionManager) IsConnected() bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.connected
}

// Reconnect reconecta a RabbitMQ
func (cm *ConnectionManager) Reconnect() error {
	cm.Close()
	
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := cm.Connect()
		if err == nil {
			return nil
		}
		
		backoff := time.Duration(i*i) * time.Second
		log.Printf("⚠️ Reconnection attempt %d/%d failed, retrying in %v", i+1, maxRetries, backoff)
		time.Sleep(backoff)
	}
	
	return fmt.Errorf("failed to reconnect after %d attempts", maxRetries)
}

// Close cierra la conexión
func (cm *ConnectionManager) Close() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.channel != nil {
		cm.channel.Close()
		cm.channel = nil
	}

	if cm.connection != nil {
		cm.connection.Close()
		cm.connection = nil
	}

	cm.connected = false
	log.Println("🔌 Disconnected from RabbitMQ")
}