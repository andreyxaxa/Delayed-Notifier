package telegrambotsender

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI struct {
	botAPI *tgbotapi.BotAPI
}

func New(token string) (*BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("TelegramBotSender - New - tgbotapi.NewBotAPI: %w", err)
	}

	return &BotAPI{botAPI: bot}, nil
}

func (b *BotAPI) Send(chatID int, message string) error {
	msg := tgbotapi.NewMessage(int64(chatID), message)

	_, err := b.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("TelegramBotSender - Send - b.botAPI.Send: %w", err)
	}

	return nil
}
