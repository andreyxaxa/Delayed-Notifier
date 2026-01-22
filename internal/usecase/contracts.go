package usecase

import (
	"context"
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/google/uuid"
)

type (
	Notification interface {
		CreateNotification(ctx context.Context, notification entity.Notification) error
		GetStatus(ctx context.Context, notificationUID uuid.UUID) (string, error)
		CancelNotification(ctx context.Context, notificationUID uuid.UUID) error
		GetPendingNotifications(ctx context.Context, now time.Time) ([]entity.Notification, error)
		MarkAsProcessing(ctx context.Context, notificationUID uuid.UUID, now time.Time) error
		RetryNotification(ctx context.Context, notificationUID uuid.UUID, nextTry, now time.Time) error
		MarkAsFailed(ctx context.Context, notificationUID uuid.UUID, now time.Time) error
	}

	// TODO: telegram, mail senders interfaces
)
