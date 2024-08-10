package main

import (
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MailPayload struct {
	To      string `json:"to"`
	Msg     string `json:"msg"`
	Subject string `json:"subject"`
}

func (p *MailPayload) toMailBody() []byte {
	msg := fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		p.To,
		p.Subject,
		p.Msg,
	)
	return []byte(msg)
}

type MailConsumer struct {
	channel *amqp.Channel
	MailHandler
}

func NewMailConsumer(channel *amqp.Channel, handler MailHandler) *MailConsumer {
	return &MailConsumer{
		channel:     channel,
		MailHandler: handler,
	}
}

func (mc *MailConsumer) startConsuming(queueName string) (msgs <-chan amqp.Delivery, err error) {
	// Creates queue if does not exists
	q, err := mc.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	msgs, err = mc.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return
}

func (mc *MailConsumer) Consume(queueName string) error {
	msgs, err := mc.startConsuming(queueName)
	if err != nil {
		return errors.New("Couldn't start consuming because: " + err.Error())
	}

	for msg := range msgs {
		go func() {
			var payload MailPayload
			err := json.Unmarshal(msg.Body, &payload)
			if err != nil {
				LOGGER.Info("Error unmarshalling mail payload: " + err.Error())
			}
			mc.sendEmail(payload)
		}()
	}

	return nil
}
