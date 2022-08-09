package application_2

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"twitch_chat_analysis/cmd/helper"
)

func Application_2() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helper.FailOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	defer conn.Close()
	helper.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	helper.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	helper.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			//	setDataToRedis(rdb, d.Body)
		}
	}()
	log.Printf("Waiting for messages")
	<-forever
}
