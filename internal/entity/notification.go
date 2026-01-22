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
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Text    string `json:"text" validate:"required"`
}

type TelegramPayload struct {
	ChatID int    `json:"chat_id" validate:"required"`
	Text   string `json:"text" validate:"required"`
}
