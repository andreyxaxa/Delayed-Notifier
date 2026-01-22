package request

import (
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
)

type CreateNotificationRequest struct {
	SendAt  time.Time      `json:"send_at" validate:"required" example:"2026-01-23T20:03:30+03:00"`
	Payload entity.Payload `json:"payload" validate:"required"`
}
