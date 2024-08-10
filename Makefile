.PHONY: mail queue user run_user

queue:
	docker compose up -d rabbitmq

mail:
	docker compose up -d mail_service

user:
	docker compose up -d user_service

all: mail user
	@echo "Running all services needed"