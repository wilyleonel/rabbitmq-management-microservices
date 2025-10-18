# 🐇 RabbitMQ Management Microservices Package

![RabbitMQ](https://www.rabbitmq.com/img/rabbitmq-logo-with-name.svg)
![Go](https://img.shields.io/badge/Go-1.19%252B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

Un paquete **robusto y flexible** para **Go** que gestiona **conexiones, producers, consumers, reintentos y dead letter queues** en **RabbitMQ**, diseñado para **arquitecturas de microservicios**.

---

## 📦 Instalación

```bash
go get github.com/wilyleonel/rabbitmq-management-microservices
```

---

## 🧩 Uso en el código

```go
import "github.com/wilyleonel/rabbitmq-management-microservices"
```

---

## ✨ Características

🔄 **Reconexión Automática** – Manejo automático de caídas de conexión
🔁 **Sistema de Reintentos** – Backoff exponencial para errores temporales
💀 **Dead Letter Queues** – Manejo elegante de mensajes fallidos
🚀 **High Performance** – Configuración optimizada para alto rendimiento
⚙️ **Configuración Flexible** – Adaptable a diferentes casos de uso
🛡️ **Manejo Robusto de Errores** – Recuperación ante fallos
🏗️ **Arquitectura Microservicios** – Diseñado para sistemas distribuidos

---

## ⚡ Configuración Básica

```go
package main

import (
    "log"
    "time"

    "github.com/wilyleonel/rabbitmq-management-microservices"
    "github.com/streadway/amqp"
)

func main() {
    config := rabbitmq_management_microservices.Config{
        URL:           "amqp://guest:guest@localhost:5672/",
        PrefetchCount: 1,
        PrefetchSize:  0,
    }

    manager, err := rabbitmq_management_microservices.NewRabbitMQManager(config)
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Close()
}
```

---

## 🏗️ Crear Exchange y Queue

```go
exchangeConfig := rabbitmq_management_microservices.ExchangeConfig{
    Name:    "notifications_exchange",
    Type:    "topic",
    Durable: true,
}
err = manager.CreateExchange(exchangeConfig)

queueConfig := rabbitmq_management_microservices.QueueConfig{
    Name:    "email_queue",
    Durable: true,
    Args:    amqp.Table{"x-dead-letter-exchange": "dlx"},
}
err = manager.CreateQueue(queueConfig)
manager.BindQueue("email_queue", "notifications_exchange", "email")
```

---

## 📤 Enviar Mensajes (Producer)

```go
producer := manager.GetProducer()

message := map[string]interface{}{
    "id": "123",
    "type": "welcome",
    "user_id": "user-456",
    "data": map[string]interface{}{"name": "John"},
}

err = producer.Publish(
    "notifications_exchange",
    "email",
    message,
    rabbitmq_management_microservices.WithMessageID("msg-123"),
    rabbitmq_management_microservices.WithHeaders(amqp.Table{"priority": "high"}),
)
```

---

## 📥 Consumir Mensajes (Consumer)

```go
consumerConfig := rabbitmq_management_microservices.ConsumerConfig{
    QueueName:   "email_queue",
    ConsumerTag: "email-consumer",
    AutoAck:     false,
}

retryConfig := rabbitmq_management_microservices.RetryConfig{
    Enabled:    true,
    MaxRetries: 3,
    Delay:      time.Second * 2,
    Backoff:    true,
}

deadLetterConfig := rabbitmq_management_microservices.DeadLetterConfig{
    Enabled:    true,
    Exchange:   "dlx",
    Queue:      "email_dlq",
    RoutingKey: "email.dead",
}

consumer := manager.CreateConsumer(
    handleEmailMessage,
    consumerConfig,
    retryConfig,
    deadLetterConfig,
)
consumer.Start()

func handleEmailMessage(msg amqp.Delivery) error {
    log.Printf("📨 Mensaje recibido: %s", string(msg.Body))
    return nil
}
```

---

## 🚀 Ejemplo Completo: Sistema de Microservicios

Incluye **Orders**, **Notifications** y **Payments** simulando un entorno distribuido con reintentos y DLQs.
*(Código completo como el mostrado anteriormente, manteniendo la estructura y funciones por servicio)*

---

## 🔧 Configuración para Microservicios

```go
retryConfig := rabbitmq_management_microservices.RetryConfig{
    Enabled:    true,
    MaxRetries: 5,
    Delay:      time.Second,
    Backoff:    true,
}

deadLetterConfig := rabbitmq_management_microservices.DeadLetterConfig{
    Enabled:    true,
    Exchange:   "dlx",
    Queue:      "microservice_dlq",
    RoutingKey: "microservice.dead",
}

consumerConfig := rabbitmq_management_microservices.ConsumerConfig{
    QueueName:   "service_queue",
    ConsumerTag: "service-consumer",
    AutoAck:     false,
    Args:        amqp.Table{"x-priority": 10},
}
```

---

## 🏗️ Patrones para Microservicios

### 1. Event Sourcing

```go
producer.Publish("orders", "order.created", orderEvent,
    rabbitmq_management_microservices.WithMessageID(orderEvent.OrderID))
```

### 2. Saga Pattern

```go
producer.Publish("payments", "payment.processed", paymentEvent,
    rabbitmq_management_microservices.WithHeaders(amqp.Table{
        "saga_id": sagaID,
        "step":    "payment",
    }))
```

### 3. CQRS

```go
producer.Publish("notifications", "user.updated", userEvent,
    rabbitmq_management_microservices.WithMessageID(userEvent.UserID))
```

---

## 📁 Estructura del Proyecto

```
rabbitmq-management-microservices/
├── connection.go        # Gestión de conexiones
├── consumer.go          # Consumers con reintentos
├── producer.go          # Producers optimizados
├── retry.go             # Sistema de reintentos
├── dead_letter.go       # Manejo de DLQ
├── manager.go           # Manager principal
├── types.go             # Tipos y configuraciones
└── examples/            # Ejemplos de microservicios
    ├── order_service/
    ├── notification_service/
    └── payment_service/
```

---

## 📋 Best Practices para Microservicios

✅ Un exchange por dominio – Separar responsabilidades
✅ Colas específicas por servicio – Aislamiento
✅ Dead letters por microservicio – Debugging fácil
✅ Tags descriptivos – Monitoreo efectivo
✅ Headers para correlación – Tracing distribuido
✅ Prefetch configurado – Balance de carga

---

## 📊 Monitoreo

```bash
rabbitmqctl list_queues name messages_ready messages_unacknowledged
rabbitmqctl list_exchanges
rabbitmqctl list_bindings
rabbitmqctl get_queue order_dlq
```

---

MIT License

Copyright (c) 2025 wilyleonel

Permission is hereby granted...