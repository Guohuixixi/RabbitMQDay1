package main

import (
	"RabbitMQDay1/util"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)

}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(n int) (res int, err error) {
	conn, err := amqp.Dial("amqp://223.26.59.136:5672/")
	util.FailOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	util.FailOnError(err, "Failed to declare queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	util.FailOnError(err, "Falied to register consumer")

	corrId := randomString(26)
	err = ch.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	util.FailOnError(err, "Failed to publish message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			util.FailOnError(err, "Failed to convert body to Integer")
			break
		}
	}
	return
}

func bodyFromFib(args []string) int {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	util.FailOnError(err, "Failed to convert arg to integer")
	return n
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	n := bodyFromFib(os.Args)

	log.Printf("[x] Request fib(%d)", n)
	res, err := fibonacciRPC(n)
	util.FailOnError(err, "Failed to handle RPC request")

	log.Printf("[.] Got %d", res)
}
