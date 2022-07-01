package adapter

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueAdapter interface {
	Publish(message string) error
}

type RabbitMQAdapter struct {
	Username string
	Password string
	Host     string
	Port     string
	Queue    string
}

func (r *RabbitMQAdapter) Publish(message string) error {
	conn, err := amqp.Dial("amqp://" + r.Username + ":" + r.Password + "@" + r.Host + ":" + r.Port)
	r.failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	r.failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.Queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	r.failOnError(err, "Failed to declare a queue")

	body := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		body)
	r.failOnError(err, "Failed to publish a message")
	return nil
}

func (r *RabbitMQAdapter) failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
