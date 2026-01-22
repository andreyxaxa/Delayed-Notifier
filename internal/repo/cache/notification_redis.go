package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	nredis "github.com/andreyxaxa/Delayed-Notifier/pkg/redis"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/errs"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type NotificationCache struct {
	c *nredis.Client
}

func New(c *nredis.Client) *NotificationCache {
	return &NotificationCache{c: c}
}

func (r *NotificationCache) GetStatus(ctx context.Context, notificationUID uuid.UUID) (string, error) {
	status, err := r.c.Client.Get(ctx, notificationUID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("NotificationCache - GetStatus - r.c.Client.Get: %w", errs.ErrRecordNotFound)
		} else {
			return "", fmt.Errorf("NotificationCache - GetStatus - r.c.Client.Get: %w", err)
		}
	}

	return status, nil
}

func (r *NotificationCache) SetStatus(ctx context.Context, notificationUID uuid.UUID, status string) error {
	// TODO: конфигурировать TTL
	if err := r.c.Client.Set(ctx, notificationUID.String(), status, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("NotificationCache - SetStatus - r.c.Client.Set: %w", err)
	}

	return nil
}

func (r *NotificationCache) DeleteStatus(ctx context.Context, notificationUID uuid.UUID) error {
	if err := r.c.Client.Del(ctx, notificationUID.String()).Err(); err != nil {
		return fmt.Errorf("NotificationCache - DeleteStatus - r.c.Client.Del: %w", err)
	}

	return nil
}
