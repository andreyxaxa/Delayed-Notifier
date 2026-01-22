package request

import (
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
)

type CreateNotificationRequest struct {
	SendAt  time.Time      `json:"send_at" validate:"required"`
	Payload entity.Payload `json:"payload" validate:"required"`
}
