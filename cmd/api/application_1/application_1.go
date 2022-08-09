package application_1

import (
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"twitch_chat_analysis/cmd/helper"
)

func Application_1() {

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

	r.Run()
}
