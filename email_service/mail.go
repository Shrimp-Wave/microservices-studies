package main

import (
	"log/slog"
	"net/smtp"
)

type MailHandler interface {
	sendEmail(payload MailPayload) error
}

type MailHandlerImpl struct {
	auth    smtp.Auth
	from    string
	service string
}

func newMailHandlerImpl(host, from, pass string) *MailHandlerImpl {
	auth := smtp.PlainAuth("", from, pass, host)
	return &MailHandlerImpl{auth: auth, from: from, service: host}
}

func (m *MailHandlerImpl) sendEmail(payload MailPayload) error {
	return smtp.SendMail(
		m.service+":587",
		m.auth,
		m.from,
		[]string{payload.To},
		payload.toMailBody(),
	)
}

type MailHandlerLogger struct {
	service MailHandler
	logger  *slog.Logger
}

func (m *MailHandlerLogger) sendEmail(payload MailPayload) (err error) {
	defer func() {
		if err != nil {
			m.logger.Error("Could not send email to target...", "TO", payload.To)
		} else {
			m.logger.Info("Email successfully sent to target", "TO", payload.To)
		}
	}()

	m.logger.Info("Trying to send email...", "TO", payload.To, "SUBJECT", payload.Subject)
	return m.service.sendEmail(payload)
}

func NewMailHandler(host, from, pass string, logger *slog.Logger) MailHandler {
	service := newMailHandlerImpl(host, from, pass)
	return &MailHandlerLogger{
		service: service,
		logger:  logger,
	}
}
