package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting listening for messages...")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("Failed to open a channel:", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"test", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Fatalln("Failed to declare a queue:", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalln("Failed to register a consumer:", err)
	}

	// Start a goroutine to process messages from RabbitMQ
	go func() {
		for msg := range msgs {
			fmt.Printf("%s Received a message: %s\n", time.Now(), msg.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// Process messages from the messageChannel in the main routine
	var forever chan struct{}
	<-forever
}
