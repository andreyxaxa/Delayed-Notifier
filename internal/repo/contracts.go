package repo

import (
	"context"
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/google/uuid"
)

type (
	NotificationRepo interface {
		CreateNotification(ctx context.Context, notification entity.Notification) error
		GetStatus(ctx context.Context, notificationUID uuid.UUID) (string, error)
		CancelNotification(ctx context.Context, notificationUID uuid.UUID) error
		GetPendingNotifications(ctx context.Context, now time.Time) ([]entity.Notification, error)
		MarkAsProcessing(ctx context.Context, notificationUID uuid.UUID, now time.Time) error
		RetryNotification(ctx context.Context, notificationUID uuid.UUID, nextTry, now time.Time) error
		MarkAsFailed(ctx context.Context, notificationUID uuid.UUID, now time.Time) error
		MarkAsSent(ctx context.Context, notificationUID uuid.UUID, now time.Time) error
	}

	CacheNotificationRepo interface {
		GetStatus(ctx context.Context, notificationUID uuid.UUID) (string, error)
		SetStatus(ctx context.Context, notificationUID uuid.UUID, status string) error
		DeleteStatus(ctx context.Context, notificationUID uuid.UUID) error
	}
)
