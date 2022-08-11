package application_1

import (
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"twitch_chat_analysis/cmd/helper"
)

var (
	r = gin.Default()
)

func Application_1() {

	//Application 1
	r.POST("/message", func(c *gin.Context) {

		//connect to RabbitMQ
		conn, err := amqp.Dial("amqp://user:password@localhost:7001/")
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

		var body map[string]string
		//convert json to map
		err = c.BindJSON(&body)
		if err != nil {
			c.JSON(400, "bad request") // Handle error
		}

		bf, _ := helper.Serialize(body)
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        bf,
			})
		helper.FailOnError(err, "Failed to publish a message")
		c.JSON(200, "Message sent successfully !!")
	})

	//Application 3
	r.GET("/message/list", func(c *gin.Context) {
		rdb := helper.ConnectToRedis()
		var body map[string]string
		//convert json to map
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(400, "bad request") // Handle error
		}
		var msg = helper.GetSortedMessagesListFromRedis(rdb, body["sender"], body["receiver"])

		helper.GetSortedMessagesListFromRedis(rdb, body["sender"], body["receiver"])

		c.JSON(200, msg)
	})

	r.Run()

}
