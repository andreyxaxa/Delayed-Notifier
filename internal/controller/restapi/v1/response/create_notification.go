package response

import "github.com/google/uuid"

type CreateNotificationResponse struct {
	UID    uuid.UUID `json:"uid"`
	Status string    `json:"status"`
}
