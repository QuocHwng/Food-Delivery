package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"notification-service/internal/model"
	"notification-service/internal/service"
)

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	svc  *service.NotificationService
}

func NewConsumer(amqpURL string, svc *service.NotificationService) (*Consumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn: conn,
		ch:   ch,
		svc:  svc,
	}, nil
}

func (c *Consumer) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *Consumer) StartConsuming() {
	// Declare an exchange (topic)
	err := c.ch.ExchangeDeclare(
		"food_delivery_events", // name
		"topic",                // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	// Declare a queue
	q, err := c.ch.QueueDeclare(
		"notification_queue", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Bind the queue to multiple routing keys
	routingKeys := []string{"order.created", "order.status_changed", "payment.processed", "payment.failed"}
	for _, rk := range routingKeys {
		err = c.ch.QueueBind(
			q.Name,                 // queue name
			rk,                     // routing key
			"food_delivery_events", // exchange
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to bind queue to %s: %s", rk, err)
		}
	}

	msgs, err := c.ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	log.Println("[RabbitMQ] 🐰 Waiting for events. To exit press CTRL+C")

	go func() {
		for d := range msgs {
			log.Printf("[RabbitMQ] Received a message: %s", d.RoutingKey)

			var event model.NotificationEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error decoding event: %s", err)
				d.Nack(false, false)
				continue
			}

			// Force type to match routing key if empty
			if event.Type == "" {
				event.Type = d.RoutingKey
			}

			if err := c.svc.HandleEvent(event); err != nil {
				log.Printf("Error handling event: %s", err)
				d.Nack(false, true) // Requeue
			} else {
				d.Ack(false)
			}
		}
	}()
}
