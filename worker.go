package main

import (
	"RabbitMQDay1/util"
	"bytes"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://223.26.59.136:5672/")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	util.FailOnError(err, "Failed to declare a queue")
	//公平分发，消费者前一个确认没有发回就不能继续为它分发任务
	err = ch.Qos(
		1,
		0,
		false,
	)
	util.FailOnError(err, "Failed to configue Qos")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, //不让接收到消息就自动返回
		false,
		false,
		false,
		nil,
	)
	util.FailOnError(err, "Failed to register a consumer")
	forever := make(chan struct{})
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false) //业务逻辑处理完成手动返回确认
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
