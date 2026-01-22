package persistent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/postgres"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/errs"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/status"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

/*
	UID       uuid.UUID `json:"uid"`
    Payload   Payload   `json:"payload"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    SendAt    time.Time `json:"send_at"`

    Status     string `json:"status"` // pendind, processing, sent, cancelled(юзер отменил), failed(провалено после N попыток)
    RetryCount int    `json:"retry_count"`
*/

const (
	// Table
	notificationsTable = "notifications"

	// Column
	uidColumn        = "uid"
	payloadColumn    = "payload"
	createdAtColumn  = "created_at"
	updatedAtColumn  = "updated_at"
	sendAtColumn     = "send_at"
	statusColumn     = "status"
	retryCountColumn = "retry_count"
)

type NotificationRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *NotificationRepo {
	return &NotificationRepo{pg}
}

func (r *NotificationRepo) CreateNotification(ctx context.Context, notification entity.Notification) error {
	payload, err := json.Marshal(notification.Payload)
	if err != nil {
		return fmt.Errorf("NotificationRepo - CreateNotification - json.Marshal: %w", err)
	}

	// Запись, при конфликте 'uid' ничего не делаем
	sql, args, err := r.Builder.
		Insert(notificationsTable).
		Columns(uidColumn, payloadColumn, createdAtColumn, updatedAtColumn, sendAtColumn, statusColumn).
		Values(notification.UID, payload, notification.CreatedAt, notification.UpdatedAt, notification.SendAt, notification.Status).
		Suffix(fmt.Sprintf("ON CONFLICT (%s) DO NOTHING", uidColumn)).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - CreateNotification - r.Builder.ToSql: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("NotificationRepo - CreateNotification - r.Pool.Exec: %w", err)
	}

	// Если не было вставки - значит был конфликт
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("NotificationRepo - CreateNotification: %w", errs.ErrRecordAlreadyExists)
	}

	return nil
}

func (r *NotificationRepo) GetStatus(ctx context.Context, notificationUID uuid.UUID) (string, error) {
	var status string

	sql, args, err := r.Builder.
		Select(statusColumn).
		From(notificationsTable).
		Where(squirrel.Eq{uidColumn: notificationUID}).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("NotificationRepo - GetStatus - r.Builder.ToSql: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("NotificationRepo - GetStatus - row.Scan: %w", errs.ErrRecordNotFound)
		}
		return "", fmt.Errorf("NotificationRepo - GetStatus - row.Scan: %w", err)
	}

	return status, nil
}

func (r *NotificationRepo) CancelNotification(ctx context.Context, notificationUID uuid.UUID) error {
	sql, args, err := r.Builder.
		Update(notificationsTable).
		Where(squirrel.Eq{uidColumn: notificationUID}).
		Set(statusColumn, status.Cancelled).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - CancelNotification - r.Builder.ToSql: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("NotificationRepo - CancelNotification - r.Pool.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("NotificationRepo - CancelNotification: %w", errs.ErrRecordNotFound)
	}

	return nil
}

func (r *NotificationRepo) GetPendingNotifications(ctx context.Context, now time.Time) ([]entity.Notification, error) {
	sql, args, err := r.Builder.
		Select(uidColumn, payloadColumn, createdAtColumn, updatedAtColumn, sendAtColumn, statusColumn, retryCountColumn).
		From(notificationsTable).
		Where(squirrel.Eq{statusColumn: status.Pending}).
		Where(squirrel.LtOrEq{sendAtColumn: now}).
		OrderBy(sendAtColumn + " ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("NotificationRepo - GetPendingNotifications: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("NotificationRepo - GetPendingNotifications - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var ntfs []entity.Notification
	for rows.Next() {
		var ntf entity.Notification

		err := rows.Scan(
			&ntf.UID,
			&ntf.Payload,
			&ntf.CreatedAt,
			&ntf.UpdatedAt,
			&ntf.SendAt,
			&ntf.Status,
			&ntf.RetryCount,
		)
		if err != nil {
			return nil, fmt.Errorf("NotificationRepo - GetPendingNotifications - rows.Scan: %w", err)
		}

		ntfs = append(ntfs, ntf)
	}

	return ntfs, nil
}

func (r *NotificationRepo) MarkAsProcessing(ctx context.Context, notificationUID uuid.UUID, now time.Time) error {
	sql, args, err := r.Builder.
		Update(notificationsTable).
		Set(statusColumn, status.Processing).
		Set(updatedAtColumn, now).
		Where(squirrel.Eq{uidColumn: notificationUID}).
		Where(squirrel.Eq{statusColumn: status.Pending}).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - MarkAsProcessing - r.Builder.ToSql: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("NotificationRepo - MarkAsProcessing - r.Pool.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("NotificationRepo - MarkAsProcessing: %w", errs.ErrRecordNotFoundOrAlreadyProcessing)
	}

	return nil
}

func (r *NotificationRepo) RetryNotification(ctx context.Context, notificationUID uuid.UUID, nextTry, now time.Time) error {
	sql, args, err := r.Builder.
		Update(notificationsTable).
		Set(statusColumn, status.Pending).
		Set(retryCountColumn, squirrel.Expr(retryCountColumn+" + 1")).
		Set(updatedAtColumn, now).
		Set(sendAtColumn, nextTry).
		Where(squirrel.Eq{uidColumn: notificationUID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - RetryNotification - r.Builder.ToSql: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("NotificationRepo - RetryNotification - r.Pool.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("NotificationRepo - RetryNotification: %w", errs.ErrRecordNotFound)
	}

	return nil
}

func (r *NotificationRepo) MarkAsFailed(ctx context.Context, notificationUID uuid.UUID, now time.Time) error {
	sql, args, err := r.Builder.
		Update(notificationsTable).
		Set(statusColumn, status.Failed).
		Set(updatedAtColumn, now).
		Where(squirrel.Eq{uidColumn: notificationUID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - MarkAsFailed - r.Builder.ToSql: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("NotificationRepo - MarkAsFailed - r.Pool.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("NotificationRepo - MarkAsFailed: %w", errs.ErrRecordNotFound)
	}

	return nil
}

func (r *NotificationRepo) MarkAsSent(ctx context.Context, notificationUID uuid.UUID, now time.Time) error {
	sql, args, err := r.Builder.
		Update(notificationsTable).
		Set(statusColumn, status.Sent).
		Set(updatedAtColumn, now).
		Where(squirrel.Eq{uidColumn: notificationUID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - MarkAsSent - r.Builder.ToSql: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("NotificationRepo - MarkAsSent - r.Pool.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("NotificationRepo - MarkAsSent: %w", errs.ErrRecordNotFound)
	}

	return nil
}
