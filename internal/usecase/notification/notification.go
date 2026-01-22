package notification

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/andreyxaxa/Delayed-Notifier/internal/repo"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/errs"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/status"
	"github.com/google/uuid"
)

type UseCase struct {
	repo  repo.NotificationRepo
	cache repo.CacheNotificationRepo

	l logger.Interface
}

func New(r repo.NotificationRepo, c repo.CacheNotificationRepo, l logger.Interface) *UseCase {
	return &UseCase{
		repo:  r,
		cache: c,
		l:     l,
	}
}

func (uc *UseCase) CreateNotification(ctx context.Context, notification entity.Notification) error {
	err := uc.repo.CreateNotification(ctx, notification)
	if err != nil {
		return fmt.Errorf("NotificationUseCase - CreateNotification - uc.repo.CreateNotification: %w", err)
	}

	return nil
}

func (uc *UseCase) GetStatus(ctx context.Context, notificationUID uuid.UUID) (string, error) {
	// searching in cache
	status, err := uc.cache.GetStatus(ctx, notificationUID)
	if err == nil {
		return status, nil
	}

	if !errors.Is(err, errs.ErrRecordNotFound) {
		uc.l.Error(err, "NotificationUseCase - GetStatus - uc.cache.GetStatus")
	}

	// searching in repo
	status, err = uc.repo.GetStatus(ctx, notificationUID)
	if err != nil {
		return "", fmt.Errorf("NotificationUseCase - GetStatus - uc.repo.GetStatus: %w", err)
	}

	// cache set
	err = uc.cache.SetStatus(ctx, notificationUID, status)
	if err != nil {
		uc.l.Warn("NotificationUseCase - GetStatus - uc.cache.SetStatus: %v", err)
	}

	return status, nil
}

func (uc *UseCase) CancelNotification(ctx context.Context, notificationUID uuid.UUID) error {
	s, err := uc.repo.GetStatus(ctx, notificationUID)
	if err != nil {
		return fmt.Errorf("NotificationUseCase - CancelNotification - uc.repo.GetStatus: %w", err)
	}

	if s == status.Cancelled {
		return fmt.Errorf("NotificationUseCase - CancelNotification: %w", errs.ErrAlreadyCancelled)
	}

	if s == status.Sent || s == status.Failed {
		return fmt.Errorf("NotificationUseCase - CancelNotification: %w", errs.ErrAlreadySentOrFailed)
	}

	err = uc.repo.CancelNotification(ctx, notificationUID)
	if err != nil {
		return fmt.Errorf("NotificationUseCase - CancelNotification - uc.repo.CancelNotification: %w", err)
	}

	// cache invalidation
	err = uc.cache.DeleteStatus(ctx, notificationUID)
	if err != nil {
		uc.l.Warn("NotificationUseCase - CancelNotification - uc.cache.DeleteStatus: %v", err)
	}

	return nil
}

func (uc *UseCase) GetPendingNotifications(ctx context.Context, now time.Time) ([]entity.Notification, error) {
	ntfs, err := uc.repo.GetPendingNotifications(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("NotificationUseCase - GetPendingNotifications - uc.repo.GetPendingNotifications: %w", err)
	}

	return ntfs, nil
}

func (uc *UseCase) MarkAsProcessing(ctx context.Context, notificationUID uuid.UUID, now time.Time) error {
	err := uc.repo.MarkAsProcessing(ctx, notificationUID, now)
	if err != nil {
		return fmt.Errorf("NotificationUseCase - MarkAsProcessing - uc.repo.MarkAsProcessing: %w", err)
	}

	// cache invalidation
	err = uc.cache.DeleteStatus(ctx, notificationUID)
	if err != nil {
		uc.l.Warn("NotificationUseCase - MarkAsProcessing - uc.cache.DeleteStatus: %v", err)
	}

	return nil
}

func (uc *UseCase) RetryNotification(ctx context.Context, notificationUID uuid.UUID, nextTry, now time.Time) error {
	err := uc.repo.RetryNotification(ctx, notificationUID, nextTry, now)
	if err != nil {
		return fmt.Errorf("NotificationUseCase - RetryNotification - uc.repo.RetryNotification: %w", err)
	}

	// cache invalidation
	err = uc.cache.DeleteStatus(ctx, notificationUID)
	if err != nil {
		uc.l.Warn("NotificationUseCase - RetryNotification - uc.cache.DeleteStatus: %v", err)
	}

	return nil
}

func (uc *UseCase) MarkAsFailed(ctx context.Context, notificationUID uuid.UUID, now time.Time) error {
	err := uc.repo.MarkAsFailed(ctx, notificationUID, now)
	if err != nil {
		return fmt.Errorf("NotificationUseCase - MarkAsFailed - uc.repo.MarkAsFailed: %w", err)
	}

	// cache invalidation
	err = uc.cache.DeleteStatus(ctx, notificationUID)
	if err != nil {
		uc.l.Warn("NotificationUseCase - MarkAsFailed - uc.cache.DeleteStatus: %v", err)
	}

	return nil
}

// TODO: telegram, mail senders
