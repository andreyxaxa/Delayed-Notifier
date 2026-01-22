package mailsender

import (
	"fmt"
	"strings"
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/smtpsender"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/errs"
)

type MailSender struct {
	sender *smtpsender.MailSender
}

func New(sender *smtpsender.MailSender) *MailSender {
	return &MailSender{
		sender: sender,
	}
}

func (s *MailSender) Send(notification entity.Notification) error {
	now := time.Now()
	if notification.SendAt.After(now) {
		return fmt.Errorf("mailsender - Send: %w", errs.ErrFutureNotification)
	}

	email := notification.Payload.Email
	toStr := strings.TrimSpace(email.To)
	text := strings.TrimSpace(email.Text)
	subject := strings.TrimSpace(email.Subject)

	to := []string{toStr}
	msg := []byte("To: " + toStr + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		text + "\r\n",
	)

	err := s.sender.SendMail(to, msg)
	if err != nil {
		return fmt.Errorf("mailsender - Send - s.sender.SendMail: %w", err)
	}

	return nil
}
