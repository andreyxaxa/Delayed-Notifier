package telegramsender

import (
	"fmt"
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/telegrambotsender"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/errs"
)

type TelegramSender struct {
	sender *telegrambotsender.BotAPI
}

func New(sender *telegrambotsender.BotAPI) *TelegramSender {
	return &TelegramSender{sender: sender}
}

func (s *TelegramSender) Send(notification entity.Notification) (bool, error) {
	now := time.Now()
	if notification.SendAt.After(now) {
		return false, fmt.Errorf("telegramsender - Send: %w", errs.ErrFutureNotification)
	}

	telegram := notification.Payload.Telegram

	err := s.sender.Send(telegram.ChatID, telegram.Text)
	if err != nil {
		return false, nil
	}

	return true, nil
}
