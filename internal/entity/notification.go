package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	UID       uuid.UUID `json:"uid"`
	Payload   Payload   `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SendAt    time.Time `json:"send_at"`

	Status     string `json:"status"`
	RetryCount int    `json:"retry_count"`
}

type Payload struct {
	Channel  string           `json:"channel" validate:"required,oneof=email telegram"`
	Email    *EmailPayload    `json:"email,omitempty"`
	Telegram *TelegramPayload `json:"telegram,omitempty"`
}

type EmailPayload struct {
	To      string `json:"to" validate:"required,email" example:"user@example.com"`
	Subject string `json:"subject" validate:"required" example:"Birtdhay"`
	Text    string `json:"text" validate:"required" example:"Hello, dont forget about Alex's birthday!"`
}

type TelegramPayload struct {
	ChatID int    `json:"chat_id" validate:"required" example:"123456789"`
	Text   string `json:"text" validate:"required" example:"Hello, here's your delayed notification."`
}
