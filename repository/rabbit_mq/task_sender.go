package rabbitMQ

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	_ "github.com/streadway/amqp"
	"http_server/domain"
	_ "log"
)

type RabbitMQSender struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

func NewRabbitMQSender(amqpURL, queueName string) (*RabbitMQSender, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ: %v", err)
		return nil, fmt.Errorf("connecting to rabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err //TODO обработать ошибки,обернуть  в ошибки такого вида
	}
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err //TODO обработать ошибки,обернуть  в ошибки такого вида
	}
	return &RabbitMQSender{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

func (r *RabbitMQSender) Send(task domain.Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		logrus.Fatalf("Failed to marshal object: %v", err)
		return err
	}

	err = r.channel.Publish(
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		logrus.Fatalf("Failed to publish message: %v", err)
		return err
	}
	logrus.Infof("Object send to RabbitMQ: %v", task)
	return nil
}

func (r *RabbitMQSender) Close() {
	r.channel.Close()
	r.connection.Close()
}
