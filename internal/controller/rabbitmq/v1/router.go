package v1

import (
	"github.com/andreyxaxa/Delayed-Notifier/internal/usecase"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/rmqserver"
)

func NewNotificationRoutes(routes map[string]rmqserver.CallHandler, n usecase.Notification, l logger.Interface) {
	r := &V1{n: n, l: l}

	{
		routes["v1.sendNotification"] = r.sendNotification()
	}
}
