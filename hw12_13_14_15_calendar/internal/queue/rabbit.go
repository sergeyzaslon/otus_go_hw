package queue

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
)

type RabbitQueue struct {
	exchangeName string
	queueName    string
	consumerTag  string
	channel      *amqp.Channel
	logger       app.Logger
}

func NewRabbitQueue(
	ctx context.Context,
	dsn string,
	exchangeName,
	queueName string,
	logger app.Logger,
) (*RabbitQueue, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ on %s: %w", dsn, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open RabbitMQ Channel on %s: %w", dsn, err)
	}

	if len(exchangeName) > 0 {
		err = ch.ExchangeDeclare(
			exchangeName,
			amqp.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare an exchanhe %s: %w", exchangeName, err)
		}
	}

	queue, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue %s: %w", queueName, err)
	}

	err = ch.QueueBind(
		queue.Name,   // name of the queue
		queue.Name,   // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	go func() {
		<-ctx.Done()
		ch.Close()
		conn.Close()
	}()

	return &RabbitQueue{
		exchangeName: exchangeName,
		queueName:    queueName,
		consumerTag:  "calendar-consumer",
		channel:      ch,
		logger:       logger,
	}, nil
}

func (q *RabbitQueue) Add(n app.Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshall notification: %w", err)
	}

	err = q.channel.Publish(
		q.exchangeName, // exchange
		q.queueName,    // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}

	return nil
}

func (q *RabbitQueue) GetNotificationChannel() (<-chan app.Notification, error) {
	ch := make(chan app.Notification)

	deliveries, err := q.channel.Consume(
		q.queueName,   // name
		q.consumerTag, // consumerTag,
		false,         // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume queue %s: %w", q.queueName, err)
	}

	go func() {
		for d := range deliveries {
			var notification app.Notification
			err := json.Unmarshal(d.Body, &notification)
			if err != nil {
				q.logger.Error("Failed to unmarshal notification message: %s", err)
				continue
			}

			ch <- notification

			d.Ack(false)
		}

		close(ch)
	}()

	return ch, nil
}
