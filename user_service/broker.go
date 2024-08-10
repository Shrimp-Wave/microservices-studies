package main

import (
	"context"
	"encoding/json"
	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"strings"
)

var email = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Obrigado por se registrar!</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            width: 80%;
            margin: auto;
            max-width: 600px;
            background: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            padding-bottom: 20px;
            border-bottom: 2px solid #007bff;
        }
        .header h1 {
            margin: 0;
            color: #007bff;
        }
        .content {
            margin: 20px 0;
        }
        .content p {
            font-size: 16px;
            line-height: 1.5;
            color: #333333;
        }
        .footer {
            text-align: center;
            padding-top: 20px;
            border-top: 1px solid #dddddd;
            font-size: 14px;
            color: #666666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Obrigado por se registrar na Shrimpwave!</h1>
        </div>
        <div class="content">
            <p>Olá %[NAME],</p>
            <p>Estamos muito felizes por você ter se registrado na Shrimpwave! Sua conta foi criada com sucesso e agora você pode aproveitar todos os nossos recursos.</p>
            <p>Se precisar de ajuda ou tiver alguma dúvida, não hesite em nos contatar.</p>
            <p>Obrigado por fazer parte da nossa comunidade!</p>
            <p>Atenciosamente,<br>A equipe da Shrimpwave</p>
        </div>
        <div class="footer">
            <p>Se você não se registrou na Shrimpwave, por favor, ignore este e-mail.</p>
        </div>
    </div>
</body>
</html>`

type MailPayload struct {
	To      string `json:"to"`
	Msg     string `json:"msg"`
	Subject string `json:"subject"`
}

func NewMailPayload(to, username, subject string) *MailPayload {
	return &MailPayload{
		To:      to,
		Msg:     strings.Replace(email, "%[NAME]", username, 1),
		Subject: subject,
	}
}

type MailBroker struct {
	queueName string
	channel   *amqp.Channel
}

func NewMailBroker(channel *amqp.Channel, queueName string) (*MailBroker, error) {
	q, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Info("Failed to declare a queue:" + err.Error())
		return nil, err
	}

	return &MailBroker{
		queueName: q.Name,
		channel:   channel,
	}, nil
}

func (mail *MailBroker) Send(payload *MailPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	body, err := json.Marshal(payload)
	if err != nil {
		log.Error("Could not marshal mail message body:" + err.Error())
		return
	}

	log.Info("Trying to send email...")
	err = mail.channel.PublishWithContext(
		ctx,
		"",
		mail.queueName,
		true,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	if err != nil {
		log.Error("Failed to publish a message:" + err.Error())
		return
	}

	log.Info("Email sent successfully!")
}
