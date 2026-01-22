package rabbitmq

import (
	v1 "github.com/andreyxaxa/Delayed-Notifier/internal/controller/rabbitmq/v1"
	"github.com/andreyxaxa/Delayed-Notifier/internal/usecase"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/rmqserver"
)

func NewRouter(n usecase.Notification, l logger.Interface) map[string]rmqserver.CallHandler {
	routes := make(map[string]rmqserver.CallHandler)

	{
		v1.NewNotificationRoutes(routes, n, l)
	}

	return routes
}
