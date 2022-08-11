package application_2

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
	"twitch_chat_analysis/cmd/helper"
)

func Application_2() {
	rdb := helper.ConnectToRedis()
	conn, err := amqp.Dial("amqp://user:password@localhost:7001/")
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
			message, _ := helper.Deserialize(d.Body)
			log.Printf("Received a message: %s", message)
			log.Printf("Message sent at: %s", d.Timestamp)
			message["date"] = time.Now().String()
			helper.SetDataToRedis(rdb, message)

		}
	}()
	log.Printf("Waiting for messages")
	//helper.GetDataFromRedis(rdb)

	<-forever

	//log.Printf("message received")

	//helper.GetDataFromRedis(rdb)
}
