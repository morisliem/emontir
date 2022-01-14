package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"os"

	"github.com/rs/zerolog/log"
)

type EmailSender interface {
	SendActivationLink(ctx context.Context, recipient string)
}

type Mailer struct {
	auth   smtp.Auth
	sender string
	server string
}

type Config struct {
	Username string
	Password string
	Sender   string
	Host     string
	Port     int
}

func NewMailtrap(m *Config) *Mailer {
	return &Mailer{
		auth:   smtp.PlainAuth("", m.Username, m.Password, m.Host),
		sender: m.Sender,
		server: fmt.Sprintf("%s:%d", m.Host, m.Port),
	}
}

func (s *Mailer) SendActivationLink(ctx context.Context, recipient, content string) {
	from := "e-montir"
	to := []string{recipient}
	link := fmt.Sprintf("%s/auth/verify?email=%s&id=%s", os.Getenv("BASE_URL"), recipient, content)
	msg := []byte(fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", recipient) +
		"Subject: Email verification\r\n\r\n" + ActivationEmailLinkTemplate(link))

	err := smtp.SendMail(s.server, s.auth, s.sender, to, msg)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}
}
