package workerpublisher

import (
	"context"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
)

// can be rabbitmq, kafka, etc...
type NotificationPublisher interface {
	PublishNotification(ctx context.Context, notification entity.Notification) error
}
