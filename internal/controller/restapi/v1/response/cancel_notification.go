package response

import "github.com/google/uuid"

type CancelNotificationResponse struct {
	UID    uuid.UUID `json:"uid"`
	Status string    `json:"status"`
}
