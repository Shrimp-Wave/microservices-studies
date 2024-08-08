.PHONY: mail queue user

queue:
	docker compose up -d rabbitmq

mail:
	docker compose up -d mail_service

user:
	docker compose up -d user_service