package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/rmqserver"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/channel"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/errs"
	"github.com/rabbitmq/amqp091-go"
)

func (r *V1) sendNotification() rmqserver.CallHandler {
	return func(ctx context.Context, d amqp091.Delivery) error {
		var ntf entity.Notification
		err := json.Unmarshal(d.Body, &ntf)
		if err != nil {
			r.l.Error(err, "rabbitmq - v1 - sendNotification")

			return fmt.Errorf("rabbitmq - v1 - sendNotification - json.Unmarshal: %w", err)
		}

		if ntf.Payload.Channel == channel.Email {
			err = r.n.SendMailNotification(ctx, ntf)
			if err != nil {
				r.l.Error(err, "rabbitmq - v1 - sendNotification")

				return fmt.Errorf("rabbitmq - v1 - sendNotification - r.n.SendMailNotification: %w", err)
			}
		} else {
			err = r.n.SendTelegramNotification(ctx, ntf)
			if err != nil {
				r.l.Error(err, "rabbitmq - v1 - sendNotification")

				if errors.Is(err, errs.ErrInvalidChatID) {
					// ack - delete from queue
					return nil
				}

				return fmt.Errorf("rabbitmq - v1 - sendNotification - r.n.SendTelegramNotification: %w", err)
			}
		}

		return nil
	}
}
