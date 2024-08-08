package main

import (
	"github.com/joho/godotenv"
	"os"
)

type Environment struct {
	MailServiceEnv
	RabbitMQEnv
}

type MailServiceEnv struct {
	MailPassword string
	FromEmail    string
	SmtpHost     string
}

type RabbitMQEnv struct {
	Queue         string
	ConnectionUrl string
}

func NewEnvironment() *Environment {
	err := godotenv.Load("../.local.env")
	if err != nil {
		LOGGER.Warn("Could not load .local.env file, using environment variables")
	} else {
		LOGGER.Info("Loaded environment variables from .local.env file")
	}

	return &Environment{
		MailServiceEnv{
			MailPassword: os.Getenv("MAIL_SERVICE_PASSWORD"),
			FromEmail:    os.Getenv("MAIL_FROM_EMAIL"),
			SmtpHost:     os.Getenv("MAIL_SMTP_HOST"),
		},
		RabbitMQEnv{
			Queue:         os.Getenv("RABBITMQ_QUEUE_NAME"),
			ConnectionUrl: os.Getenv("RABBITMQ_URL"),
		},
	}
}
