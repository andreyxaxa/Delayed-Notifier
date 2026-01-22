package response

import "github.com/google/uuid"

type CancelNotificationResponse struct {
	UID    uuid.UUID `json:"uid" example:"9a88d642-6c65-4f0f-b8f0-b920182cceb3"`
	Status string    `json:"status" example:"pending"`
}
