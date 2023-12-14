package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func FailOnError1(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func main() {
	conn, err := amqp.Dial("amqp://223.26.59.136:5672/")
	FailOnError1(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError1(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("hello",
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError1(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError1(err, "Failed to register a consumer")
	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
