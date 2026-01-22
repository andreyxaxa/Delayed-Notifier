package smtpsender

import (
	"fmt"
	"net"
	"net/smtp"
)

const (
	_defaultSMTPHost = "smtp.gmail.com"
	_defaultSMTPPort = "587"
)

type MailSender struct {
	Username string
	Password string
	host     string
	port     string
	Addr     string
	auth     smtp.Auth
}

func New(opts ...Option) *MailSender {
	s := &MailSender{
		host: _defaultSMTPHost,
		port: _defaultSMTPPort,
		Addr: net.JoinHostPort(_defaultSMTPHost, _defaultSMTPPort),
	}

	for _, opt := range opts {
		opt(s)
	}

	auth := smtp.PlainAuth("", s.Username, s.Password, s.host)

	s.auth = auth
	s.Addr = net.JoinHostPort(s.host, s.port)

	return s
}

func (s *MailSender) SendMail(to []string, msg []byte) error {
	err := smtp.SendMail(s.Addr, s.auth, s.Username, to, msg)
	if err != nil {
		return fmt.Errorf("SMTPSender - SendMail - smtp.SendMail: %w", err)
	}

	return nil
}
