package rabbitmq

import (
	"fmt"
	"log"
	"tj/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

type RabbitConfig struct {
	ExchangeName string
	ExchangeType string
	QueueName    string
	RoutingKey   string
	ConsumerName string
}

func Connect() (*RabbitClient, error) {
	conn, err := amqp.Dial(config.Cfg.RabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	log.Println("RabbitMQ connected")
	return &RabbitClient{Conn: conn, Channel: ch}, nil
}

func (rmq *RabbitClient) Close() {
	if rmq.Channel != nil {
		_ = rmq.Channel.Close()
	}
	if rmq.Conn != nil {
		_ = rmq.Conn.Close()
	}
}

func SetupRMQ(rmq *RabbitClient, cfg RabbitConfig) error {
	exType := cfg.ExchangeType
	if exType == "" {
		exType = "topic"
	}
	if err := rmq.Channel.ExchangeDeclare(
		cfg.ExchangeName,
		exType,
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare error: %w", err)
	}

	q, err := rmq.Channel.QueueDeclare(
		cfg.QueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare error: %w", err)
	}
	if err := rmq.Channel.QueueBind(
		q.Name,
		cfg.RoutingKey,
		cfg.ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("queue bind error: %w", err)
	}

	log.Printf("RabbitMQ setup done: exchange=%s queue=%s routing=%s",
		cfg.ExchangeName, cfg.QueueName, cfg.RoutingKey)

	return nil
}

func PublishRMQ(rmq *RabbitClient, exchange, routingKey string, payload []byte) error {
	return rmq.Channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}

func ConsumeRMQWithConfig(rmq *RabbitClient, cfg RabbitConfig, autoAck bool) (<-chan amqp.Delivery, error) {
	consumer := cfg.ConsumerName
	if consumer == "" {
		consumer = "" // biar RabbitMQ generate random consumer tag
	}

	return rmq.Channel.Consume(
		cfg.QueueName,
		consumer,
		autoAck,
		false,
		false,
		false,
		nil,
	)
}
