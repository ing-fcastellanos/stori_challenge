package utils

import (
	"os"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	DialAndSend(message gomail.Message) error
}

type SMTPMailer struct {
	dialer *gomail.Dialer
}

func NewSMTPMailer() *SMTPMailer {
	var SMTPHost = os.Getenv("SMTPHost")
	var SMTPUser = os.Getenv("SMTPUser")
	var SMTPPassword = os.Getenv("SMTPPassword")
	return &SMTPMailer{
		dialer: gomail.NewDialer(SMTPHost, 587, SMTPUser, SMTPPassword),
	}
}

func (s *SMTPMailer) DialAndSend(message gomail.Message) error {
	return s.dialer.DialAndSend(&message)
}
