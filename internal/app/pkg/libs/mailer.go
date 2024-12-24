package libs

import (
	"fmt"
	"net/smtp"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
)

type MailerService interface {
	SendEmail(to []string, subject, body string) error
	SendHTMLEmail(to []string, subject, htmlBody string) error
}

type mailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type mailerServiceImpl struct {
	config mailerConfig
	logger logs.Logger
}

func NewMailerService(host string, port int, username, password, from string, logger logs.Logger) MailerService {
	return &mailerServiceImpl{
		config: mailerConfig{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			From:     from,
		},
		logger: logger,
	}
}

func (m *mailerServiceImpl) SendEmail(to []string, subject, body string) error {
	m.logger.Info("Attempting to send email")

	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", m.config.From, to[0], subject, body))

	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	err := smtp.SendMail(addr, auth, m.config.From, to, msg)
	if err != nil {
		m.logger.Error("Failed to send email", err)
		return err
	}

	m.logger.Info("Email sent successfully")
	return nil
}

func (m *mailerServiceImpl) SendHTMLEmail(to []string, subject, htmlBody string) error {
	m.logger.Info("Attempting to send HTML email")

	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
		"\r\n"+
		"%s\r\n", m.config.From, to[0], subject, htmlBody))

	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	err := smtp.SendMail(addr, auth, m.config.From, to, msg)
	if err != nil {
		m.logger.Error("Failed to send HTML email", err)
		return err
	}

	m.logger.Info("HTML email sent successfully")
	return nil
}
