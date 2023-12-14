package main

import (
	"RabbitMQDay1/util"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
)

func fib(n int) int {
	if n <= 1 {
		return n
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://223.26.59.136:5672/")
	util.FailOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	util.FailOnError(err, "Failed to declare queue")

	err = ch.Qos(
		1,
		0,
		false,
	)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			util.FailOnError(err, "Failed to convert Body to Integer")
			log.Printf(" [.]fib(%d)", n)
			response := fib(n)

			err = ch.Publish(
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(response)),
				})
			util.FailOnError(err, "Failed to publish a message")
			d.Ack(false)
		}
	}()
	log.Printf("[*] Waiting for RPC request")
	<-forever
}
