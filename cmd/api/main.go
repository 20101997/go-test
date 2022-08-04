package main

import (
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var (
	rdb = connectToRedis()
)

func main() {
	r := gin.Default()
	//Application 1
	r.POST("/message", func(c *gin.Context) {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Panicf("%s:", err)
		}
		ch, err := conn.Channel()
		if err != nil {
			log.Panicf("%s:", err)
		}
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"hello", // name
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Panicf("%s:", err)
		}
		var body message
		err = c.BindJSON(&body)
		if err != nil {
			c.JSON(400, "bad request") // Handle error
		}
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Panicf("%s:", err)
		}
		c.JSON(200, "worked")
	})

	application_2()

	//Application 3
	r.GET("/message/list", func(c *gin.Context) {
		rdb.getDataFromRedis()
		c.JSON(200, "worked")
	})

	r.Run()
}

func application_2() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("%s:", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s:", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panicf("%s:", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicf("%s:", err)
	}

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
