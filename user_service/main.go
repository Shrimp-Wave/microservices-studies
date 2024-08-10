package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	requestIdHeader = "X-Request-ID"
)

// TODO: Implement user creation send email via email service
// TODO: Implement user buy subscription calling RUST buying service
// TODO: Implement Python service to validate if transaction can be done successfully
// TODO: Implement database in everything

func stablishConnection(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:" + err.Error())
	}
	return conn
}

func getChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	return ch
}

func main() {
	// Startup Components
	e := echo.New()

	// Logger
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.Info("New Echo Instance oppenning...")

	// Validation
	e.Validator = NewCustomValidator()

	// Environment variables
	env := NewEnvironment()

	// Defining Middlewares
	e.Use(RequestIdMiddleware)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// "Database"
	persistence := NewInMemoryPersistenceLayer()

	// Message Broken to send messages
	e.Logger.Info("Trying to connect to RabbitMQ...")
	conn := stablishConnection(env.ConnectionUrl)
	defer conn.Close()

	e.Logger.Info("Getting channel from RabbitMQ connection...")
	ch := getChannel(conn)
	defer ch.Close()

	e.Logger.Info("Trying to create mail mail...")
	mail, err := NewMailBroker(ch, env.RabbitMQEnv.Queue)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	// Service instance to be used as a handler for echo methods
	service := NewUserService(persistence, mail)

	// Defining Routes
	e.POST("/", service.CreateUser, CreateUserValidatorMiddleware)
	e.GET("/:id", service.GetUserById)

	e.Logger.Fatal(e.Start(env.ServicePort))
}
