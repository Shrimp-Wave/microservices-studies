package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"log/slog"
	"os"
)

var LOGGER *slog.Logger

func stablishConnection(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ:" + err.Error())
	}
	return conn
}

func getChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("Failed to open a channel:", err)
	}
	return ch
}

func main() {
	LOGGER = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	env := NewEnvironment()

	LOGGER.Info("Trying to connect to RabbitMQ...")
	conn := stablishConnection(env.ConnectionUrl)
	defer conn.Close()

	LOGGER.Info("Getting channel from RabbitMQ connection...")
	ch := getChannel(conn)
	defer ch.Close()

	LOGGER.Info("Creating the mailing service handler...")
	handler := NewMailHandler(
		env.SmtpHost,
		env.FromEmail,
		env.MailPassword,
		LOGGER,
	)

	LOGGER.Info("Starting listening for messages...")
	err := NewMailConsumer(ch, handler).Consume(env.Queue)
	if err != nil {
		log.Fatalln("Failed to consume from RabbitMQ:", err)
	}
}
