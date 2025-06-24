package rabbitMQ

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	_ "github.com/streadway/amqp"
	"http_server/domain"
	"log"
	_ "log"
)

type ObjectSender interface {
	Send(object domain.Task)
}

type RabbitMQSender struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

func NewRabbitMQSender(amqpURL, queueName string) (*RabbitMQSender, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("connecting to rabbitMQ: %w", err) //TODO обработать ошибки,обернуть  в ошибки такого вида
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
		return err
	}

	log.Println("Object send to RabbitMQ:", task)
	return nil
}

func (r *RabbitMQSender) Close() {
	r.channel.Close()
	r.connection.Close()
}
