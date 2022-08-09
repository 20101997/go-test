package main

import (
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"twitch_chat_analysis/cmd/helper"
	"twitch_chat_analysis/cmd/model"
)

var (
	rdb  = connectToRedis()
	body model.Message
)

func main() {
	r := gin.Default()
	//Application 1
	r.POST("/message", func(c *gin.Context) {

		//connect to RabbitMQ
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		helper.FailOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		//create a channel
		ch, err := conn.Channel()
		helper.FailOnError(err, "Failed to open a channel")
		defer ch.Close()

		//declare a queue
		q, err := ch.QueueDeclare(
			"hello", // name
			false,
			false,
			false,
			false,
			nil,
		)
		helper.FailOnError(err, "Failed to declare a queue")

		/*	err = c.BindJSON(&body)
			if err != nil {
				c.JSON(400, "bad request") // Handle error
			}*/
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("hello"),
			})
		helper.FailOnError(err, "Failed to publish a message")
		c.JSON(200, "worked")
	})

	application_2()

	/*	//Application 3
		r.GET("/message/list", func(c *gin.Context) {
			rdb.getDataFromRedis()
			c.JSON(200, "worked")
		})

		r.Run()*/
}

func application_2() {

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
			setDataToRedis(rdb, d.Body)
		}
	}()
	log.Printf("Waiting for messages")
	<-forever
}
