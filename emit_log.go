package main

import (
	"RabbitMQDay1/util"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func main() {
	conn, err := amqp.Dial("amqp://223.26.59.136:5672/")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Falied to open channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	util.FailOnError(err, "Failed to declare an exchange")

	body := util.BodyFrom(os.Args)
	err = ch.Publish(
		"logs",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	util.FailOnError(err, "Failed to Publish a message")
	log.Printf("[x] Sent %s", body)

}
