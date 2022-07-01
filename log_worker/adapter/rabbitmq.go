package adapter

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber interface {
	Subscribe(callback func(<-chan amqp.Delivery, <-chan time.Time) error, tick <-chan time.Time)
}

type RabbitMQAdapter struct {
	Username string
	Password string
	Host     string
	Port     string
	Queue    string
}

func (r *RabbitMQAdapter) failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (r *RabbitMQAdapter) Subscribe(callback func(<-chan amqp.Delivery, <-chan time.Time) error, tick <-chan time.Time) {
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

	err = ch.Qos(
		5,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	r.failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	r.failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go callback(msgs, tick)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}
