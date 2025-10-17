RabbitMQ Management Microservices Package
https://www.rabbitmq.com/img/rabbitmq-logo-with-name.svg
https://img.shields.io/badge/Go-1.19%252B-blue.svg
https://img.shields.io/badge/License-MIT-green.svg

Un paquete robusto y flexible para Go que gestiona conexiones, producers, consumers, reintentos y dead letter queues en RabbitMQ para arquitecturas de microservicios.

Instalación 📦
bash
go get github.com/wilyleonel/rabbitmq-management-microservices
Uso en el código:
go
import "github.com/wilyleonel/rabbitmq-management-microservices"
Características ✨
🔄 Reconexión Automática - Manejo automático de caídas de conexión

🔁 Sistema de Reintentos - Backoff exponencial para errores temporales

💀 Dead Letter Queues - Manejo elegante de mensajes fallidos

🚀 High Performance - Configuración optimizada para alto rendimiento

⚙️ Configuración Flexible - Adaptable a diferentes casos de uso

🛡️ Manejo Robusto de Errores - Recuperación ante fallos

🏗️ Arquitectura Microservicios - Diseñado para sistemas distribuidos

Configuración Básica ⚡
go
package main

import (
"log"
"time"

    "github.com/wilyleonel/rabbitmq-management-microservices"
    "github.com/streadway/amqp"

)

func main() {
// Configuración de RabbitMQ
config := rabbitmq-management-microservices.Config{
URL: "amqp://guest:guest@localhost:5672/",
PrefetchCount: 1,
PrefetchSize: 0,
}

    // Crear manager
    manager, err := rabbitmq-management-microservices.NewRabbitMQManager(config)
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Close()

}
Crear Exchange y Queue 🏗️
go
// Crear exchange
exchangeConfig := rabbitmq-management-microservices.ExchangeConfig{
Name: "notifications_exchange",
Type: "topic",
Durable: true,
}
err = manager.CreateExchange(exchangeConfig)
if err != nil {
log.Fatal(err)
}

// Crear cola con DLQ configurada
queueConfig := rabbitmq-management-microservices.QueueConfig{
Name: "email_queue",
Durable: true,
Args: amqp.Table{
"x-dead-letter-exchange": "dlx",
},
}
err = manager.CreateQueue(queueConfig)
if err != nil {
log.Fatal(err)
}

// Bindear cola al exchange
err = manager.BindQueue("email_queue", "notifications_exchange", "email")
if err != nil {
log.Fatal(err)
}
Enviar Mensajes (Producer) 📤
go
// Obtener producer
producer := manager.GetProducer()

// Mensaje simple
message := map[string]interface{}{
"id": "123",
"type": "welcome",
"user_id": "user-456",
"data": map[string]interface{}{"name": "John"},
}

// Enviar mensaje con opciones
err = producer.Publish(
"notifications_exchange",
"email",
message,
rabbitmq-management-microservices.WithMessageID("msg-123"),
rabbitmq-management-microservices.WithHeaders(amqp.Table{"priority": "high"}),
)
if err != nil {
log.Fatal(err)
}
Consumir Mensajes (Consumer) 📥
go
// Configurar consumer
consumerConfig := rabbitmq-management-microservices.ConsumerConfig{
QueueName: "email_queue",
ConsumerTag: "email-consumer",
AutoAck: false,
}

// Configurar reintentos
retryConfig := rabbitmq-management-microservices.RetryConfig{
Enabled: true,
MaxRetries: 3,
Delay: time.Second \* 2,
Backoff: true,
}

// Configurar dead letter
deadLetterConfig := rabbitmq-management-microservices.DeadLetterConfig{
Enabled: true,
Exchange: "dlx",
Queue: "email_dlq",
RoutingKey: "email.dead",
}

// Crear consumer
consumer := manager.CreateConsumer(
handleEmailMessage,
consumerConfig,
retryConfig,
deadLetterConfig,
)

// Iniciar consumer
err = consumer.Start()
if err != nil {
log.Fatal(err)
}

// Manejar mensajes
func handleEmailMessage(msg amqp.Delivery) error {
log.Printf("📨 Mensaje recibido: %s", string(msg.Body))

    // Procesar mensaje aquí
    // Si retornas error, se activa el sistema de reintentos/DLQ
    // Si retornas nil, el mensaje se confirma (ACK)

    return nil // Éxito - mensaje será confirmado

}
Ejemplo Completo: Sistema de Microservicios 🚀
go
package main

import (
"encoding/json"
"fmt"
"log"
"time"

    "github.com/wilyleonel/rabbitmq-management-microservices"
    "github.com/streadway/amqp"

)

// Estructuras para diferentes microservicios
type OrderEvent struct {
OrderID string `json:"order_id"`
UserID string `json:"user_id"`
Amount float64 `json:"amount"`
Status string `json:"status"`
Timestamp time.Time `json:"timestamp"`
}

type NotificationEvent struct {
ID string `json:"id"`
Type string `json:"type"`
UserID string `json:"user_id"`
Data map[string]interface{} `json:"data"`
Timestamp time.Time `json:"timestamp"`
}

type PaymentEvent struct {
PaymentID string `json:"payment_id"`
OrderID string `json:"order_id"`
Amount float64 `json:"amount"`
Status string `json:"status"`
Timestamp time.Time `json:"timestamp"`
}

func main() {
// Configuración
config := rabbitmq-management-microservices.Config{
URL: "amqp://guest:guest@localhost:5672/",
PrefetchCount: 1,
}

    manager, err := rabbitmq-management-microservices.NewRabbitMQManager(config)
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Close()

    // Setup de infraestructura para microservicios
    setupMicroservicesInfrastructure(manager)

    // Iniciar microservicios
    go startOrderService(manager)
    go startNotificationService(manager)
    go startPaymentService(manager)

    // Mantener el programa corriendo
    log.Println("🚀 Ecosistema de microservicios iniciado")
    select {}

}

func setupMicroservicesInfrastructure(manager \*rabbitmq-management-microservices.RabbitMQManager) {
// Exchanges para diferentes dominios
exchanges := []rabbitmq-management-microservices.ExchangeConfig{
{Name: "orders", Type: "topic", Durable: true},
{Name: "notifications", Type: "topic", Durable: true},
{Name: "payments", Type: "topic", Durable: true},
{Name: "dlx", Type: "topic", Durable: true},
}

    for _, exchange := range exchanges {
        manager.CreateExchange(exchange)
    }

    // Colas para cada microservicio
    queues := []struct {
        name  string
        topic string
    }{
        {"order_created_queue", "order.created"},
        {"order_processed_queue", "order.processed"},
        {"email_notifications_queue", "notification.email"},
        {"push_notifications_queue", "notification.push"},
        {"payment_processed_queue", "payment.processed"},
    }

    for _, queue := range queues {
        manager.CreateQueue(rabbitmq-management-microservices.QueueConfig{
            Name:    queue.name,
            Durable: true,
            Args:    amqp.Table{"x-dead-letter-exchange": "dlx"},
        })

        // Bind basado en el tipo de exchange
        var exchange string
        switch {
        case queue.name == "order_created_queue" || queue.name == "order_processed_queue":
            exchange = "orders"
        case queue.name == "email_notifications_queue" || queue.name == "push_notifications_queue":
            exchange = "notifications"
        case queue.name == "payment_processed_queue":
            exchange = "payments"
        }

        manager.BindQueue(queue.name, exchange, queue.topic)
    }

}

// Microservicio de Órdenes
func startOrderService(manager \*rabbitmq-management-microservices.RabbitMQManager) {
// Consumer para órdenes creadas
orderConsumer := manager.CreateConsumer(
handleOrderCreated,
rabbitmq-management-microservices.ConsumerConfig{
QueueName: "order_created_queue",
AutoAck: false,
},
rabbitmq-management-microservices.RetryConfig{
Enabled: true,
MaxRetries: 3,
Delay: time.Second,
Backoff: true,
},
rabbitmq-management-microservices.DeadLetterConfig{
Enabled: true,
Exchange: "dlx",
Queue: "order_dlq",
RoutingKey: "order.dead",
},
)

    // Producer para publicar eventos
    producer := manager.GetProducer()

    // Simular creación de órdenes
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()

        orderCount := 0
        for range ticker.C {
            order := OrderEvent{
                OrderID:   fmt.Sprintf("order-%d", orderCount),
                UserID:    fmt.Sprintf("user-%d", orderCount%10),
                Amount:    99.99 + float64(orderCount),
                Status:    "created",
                Timestamp: time.Now(),
            }

            err := producer.Publish(
                "orders",
                "order.created",
                order,
                rabbitmq-management-microservices.WithMessageID(order.OrderID),
            )

            if err != nil {
                log.Printf("❌ Error publicando orden: %v", err)
            } else {
                log.Printf("✅ Orden publicada: %s", order.OrderID)
            }
            orderCount++
        }
    }()

    orderConsumer.Start()

}

func handleOrderCreated(msg amqp.Delivery) error {
var order OrderEvent
if err := json.Unmarshal(msg.Body, &order); err != nil {
return fmt.Errorf("error decodificando orden: %w", err)
}

    log.Printf("🛒 Procesando orden: %s para usuario: %s", order.OrderID, order.UserID)

    // Simular procesamiento de orden
    time.Sleep(2 * time.Second)

    // Publicar notificación
    producer := rabbitmq-management-microservices.GetProducer()
    notification := NotificationEvent{
        ID:        fmt.Sprintf("notif-%s", order.OrderID),
        Type:      "order_created",
        UserID:    order.UserID,
        Data:      map[string]interface{}{"order_id": order.OrderID, "amount": order.Amount},
        Timestamp: time.Now(),
    }

    err := producer.Publish(
        "notifications",
        "notification.email",
        notification,
        rabbitmq-management-microservices.WithMessageID(notification.ID),
    )

    if err != nil {
        return fmt.Errorf("error publicando notificación: %w", err)
    }

    log.Printf("✅ Orden procesada: %s", order.OrderID)
    return nil

}

// Microservicio de Notificaciones
func startNotificationService(manager _rabbitmq-management-microservices.RabbitMQManager) {
// Consumer para notificaciones email
emailConsumer := manager.CreateConsumer(
handleEmailNotification,
rabbitmq-management-microservices.ConsumerConfig{
QueueName: "email_notifications_queue",
AutoAck: false,
},
rabbitmq-management-microservices.RetryConfig{
Enabled: true,
MaxRetries: 3,
Delay: time.Second _ 2,
Backoff: true,
},
rabbitmq-management-microservices.DeadLetterConfig{
Enabled: true,
Exchange: "dlx",
Queue: "notification_dlq",
RoutingKey: "notification.dead",
},
)

    // Consumer para notificaciones push
    pushConsumer := manager.CreateConsumer(
        handlePushNotification,
        rabbitmq-management-microservices.ConsumerConfig{
            QueueName: "push_notifications_queue",
            AutoAck:   false,
        },
        rabbitmq-management-microservices.RetryConfig{
            Enabled:   true,
            MaxRetries: 5,
            Delay:     time.Second,
            Backoff:   true,
        },
        rabbitmq-management-microservices.DeadLetterConfig{
            Enabled:    true,
            Exchange:   "dlx",
            Queue:      "notification_dlq",
            RoutingKey: "notification.dead",
        },
    )

    go emailConsumer.Start()
    pushConsumer.Start()

}

func handleEmailNotification(msg amqp.Delivery) error {
var notification NotificationEvent
if err := json.Unmarshal(msg.Body, &notification); err != nil {
return fmt.Errorf("error decodificando notificación: %w", err)
}

    log.Printf("📧 Enviando email a usuario: %s - Tipo: %s", notification.UserID, notification.Type)
    time.Sleep(1 * time.Second)
    log.Printf("✅ Email enviado: %s", notification.ID)
    return nil

}

func handlePushNotification(msg amqp.Delivery) error {
var notification NotificationEvent
if err := json.Unmarshal(msg.Body, &notification); err != nil {
return err
}

    log.Printf("📱 Enviando push a usuario: %s - Tipo: %s", notification.UserID, notification.Type)
    time.Sleep(500 * time.Millisecond)
    log.Printf("✅ Push enviado: %s", notification.ID)
    return nil

}

// Microservicio de Pagos
func startPaymentService(manager _rabbitmq-management-microservices.RabbitMQManager) {
paymentConsumer := manager.CreateConsumer(
handlePayment,
rabbitmq-management-microservices.ConsumerConfig{
QueueName: "payment_processed_queue",
AutoAck: false,
},
rabbitmq-management-microservices.RetryConfig{
Enabled: true,
MaxRetries: 5,
Delay: time.Second _ 3,
Backoff: true,
},
rabbitmq-management-microservices.DeadLetterConfig{
Enabled: true,
Exchange: "dlx",
Queue: "payment_dlq",
RoutingKey: "payment.dead",
},
)

    paymentConsumer.Start()

}

func handlePayment(msg amqp.Delivery) error {
var payment PaymentEvent
if err := json.Unmarshal(msg.Body, &payment); err != nil {
return fmt.Errorf("error decodificando pago: %w", err)
}

    log.Printf("💳 Procesando pago: %s para orden: %s", payment.PaymentID, payment.OrderID)
    time.Sleep(3 * time.Second)
    log.Printf("✅ Pago procesado: %s", payment.PaymentID)
    return nil

}
Configuración para Microservicios 🔧
go
// Configuración de Reintentos
retryConfig := rabbitmq-management-microservices.RetryConfig{
Enabled: true,
MaxRetries: 5,
Delay: time.Second,
Backoff: true,
}

// Configuración de Dead Letter
deadLetterConfig := rabbitmq-management-microservices.DeadLetterConfig{
Enabled: true,
Exchange: "dlx",
Queue: "microservice_dlq",
RoutingKey: "microservice.dead",
}

// Configuración de Consumer
consumerConfig := rabbitmq-management-microservices.ConsumerConfig{
QueueName: "service_queue",
ConsumerTag: "service-consumer",
AutoAck: false,
Args: amqp.Table{
"x-priority": 10,
},
}
Patrones para Microservicios 🏗️

1. Event Sourcing
   go
   // Publicar evento de dominio
   err = producer.Publish(
   "orders",
   "order.created",
   orderEvent,
   rabbitmq-management-microservices.WithMessageID(orderEvent.OrderID),
   )
2. Saga Pattern
   go
   // Coordinar transacciones distribuidas
   err = producer.Publish(
   "payments",
   "payment.processed",
   paymentEvent,
   rabbitmq-management-microservices.WithHeaders(amqp.Table{
   "saga_id": sagaID,
   "step": "payment",
   }),
   )
3. CQRS
   go
   // Separar lecturas y escrituras
   err = producer.Publish(
   "notifications",
   "user.updated",
   userEvent,
   rabbitmq-management-microservices.WithMessageID(userEvent.UserID),
   )
   Estructura del Proyecto 📁
   text
   rabbitmq-management-microservices/
   ├── connection.go # Gestión de conexiones
   ├── consumer.go # Consumers con reintentos
   ├── producer.go # Producers optimizados
   ├── retry.go # Sistema de reintentos
   ├── dead_letter.go # Manejo de DLQ
   ├── manager.go # Manager principal
   ├── types.go # Tipos y configuraciones
   └── examples/ # Ejemplos de microservicios
   ├── order_service/
   ├── notification_service/
   └── payment_service/
   Best Practices para Microservicios 📋
   ✅ Un exchange por dominio - Separar concerns

✅ Colas específicas por servicio - Aislamiento

✅ Dead letters por microservicio - Debugging fácil

✅ Tags descriptivos - Monitoreo efectivo

✅ Headers para correlación - Tracing distribuido

✅ Prefetch configurado - Balance de carga

Monitoreo 📊
bash

# Ver estado de microservicios

rabbitmqctl list_queues name messages_ready messages_unacknowledged

# Ver exchanges

rabbitmqctl list_exchanges

# Ver bindings

rabbitmqctl list_bindings

# Monitorear DLQs

rabbitmqctl get_queue order_dlq
Licencia 📄
MIT License - ver el archivo LICENSE para detalles
