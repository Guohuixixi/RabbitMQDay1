package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"

	"RabbitMQDay1/util"
)

func main() {
	conn, err := amqp.Dial("amqp://223.26.59.136:5672/")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue",
		true, //更改队列持久化同时应用于生产者和消费者代码
		false,
		false,
		false,
		nil,
	)
	util.FailOnError(err, "Failed to declare a queue")
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	body := util.BodyFrom(os.Args)
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //消息持久化
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	util.FailOnError(err, "Failed to publish a message")
	log.Printf("[x] sent %s\n", body)

}
