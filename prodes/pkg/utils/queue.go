package utils

import (
	"fmt"
	"log"
	"prodes/pkg/config"

	"github.com/streadway/amqp"
)

// RabbitMQConfig contiene los datos de configuración para RabbitMQ
type RabbitMQConfig struct {
    URL      string
    Queue    string
    Exchange string
    Key      string
}

// SetupRabbitMQ establece la conexión con RabbitMQ y declara la cola
func SetupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
    conn, err := amqp.Dial(config.RabbitMQURL)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, nil, fmt.Errorf("failed to open a channel: %v", err)
    }

    // Declarar la cola
    _, err = ch.QueueDeclare(
        config.RabbitMQQueue, // Nombre de la cola
        true,                 // Duradera
        false,                // No autodelete
        false,                // No exclusiva
        false,                // No espera confirmación
        nil,                  // Argumentos adicionales
    )
    if err != nil {
        return nil, nil, fmt.Errorf("failed to declare queue: %v", err)
    }

    return conn, ch, nil
}

// ConsumeMessages consume los mensajes de la cola configurada
func ConsumeMessages(ch *amqp.Channel) (<-chan amqp.Delivery, error) {
    msgs, err := ch.Consume(
        config.RabbitMQQueue, // Nombre de la cola
        "",                   // Consumidor
        true,                 // Auto-Acknowledge
        false,                // No exclusiva
        false,                // No espera de confirmación
        false,                // Sin argumentos adicionales
        nil,                  // Argumentos adicionales
    )
    if err != nil {
        return nil, fmt.Errorf("failed to register a consumer: %v", err)
    }

    return msgs, nil
}

// HandleMessage es una función para procesar los mensajes que llegan
func HandleMessage(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        log.Printf("Recibido un mensaje: %s", d.Body)
        // Aquí puedes agregar lógica para procesar el mensaje recibido
    }
}

// CloseRabbitMQ cierra la conexión y el canal de RabbitMQ
func CloseRabbitMQ(conn *amqp.Connection, ch *amqp.Channel) {
    ch.Close()
    conn.Close()
}
