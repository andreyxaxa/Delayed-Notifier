package rmqpublisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/wb-go/wbf/rabbitmq"
)

type NotificationPublisher struct {
	pub *rabbitmq.Publisher

	exchange    string
	contentType string
	routingKey  string
}

const (
	_defaultContentType  = "application/json"
	_defaultExchangeName = "notifications"
	_defaultRoutingKey   = "v1.sendNotification"
)

func New(client *rabbitmq.RabbitClient, opts ...Option) (*NotificationPublisher, error) {
	np := &NotificationPublisher{
		pub:         nil,
		exchange:    _defaultExchangeName,
		contentType: _defaultContentType,
		routingKey:  _defaultRoutingKey,
	}

	for _, opt := range opts {
		opt(np)
	}

	rp := rabbitmq.NewPublisher(client, np.exchange, np.contentType)

	np.pub = rp

	return np, nil
}

func (p *NotificationPublisher) PublishNotification(ctx context.Context, notification entity.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("rabbitmq - publisher - PublishNotification - json.Marshal: %w", err)
	}

	err = p.pub.Publish(
		ctx,
		body,
		p.routingKey,
	)
	if err != nil {
		return fmt.Errorf("rabbitmq - publisher - PublishNotification - p.pub.Publish: %w", err)
	}

	return nil
}
