package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Environment struct {
	UserServiceEnv
	RabbitMQEnv
}

type UserServiceEnv struct {
	ServicePort string
}

type RabbitMQEnv struct {
	Queue         string
	ConnectionUrl string
}

func NewEnvironment() *Environment {
	err := godotenv.Load("../.local.env")
	if err != nil {
		slog.Warn("Could not load .local.env file, using environment variables")
	} else {
		slog.Info("Loaded environment variables from .local.env file")
	}

	return &Environment{
		UserServiceEnv{
			ServicePort: os.Getenv("USER_SERVICE_PORT"),
		},
		RabbitMQEnv{
			Queue:         os.Getenv("RABBITMQ_QUEUE_NAME"),
			ConnectionUrl: os.Getenv("RABBITMQ_URL"),
		},
	}
}
