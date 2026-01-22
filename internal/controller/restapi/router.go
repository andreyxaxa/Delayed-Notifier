package restapi

import (
	v1 "github.com/andreyxaxa/Delayed-Notifier/internal/controller/restapi/v1"
	"github.com/andreyxaxa/Delayed-Notifier/internal/usecase"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// TODO: в будущем принимать конфиг
func NewRouter(app *fiber.App, n usecase.Notification, l logger.Interface) {
	apiV1Group := app.Group("/v1")
	{
		v1.NewNotificationRoutes(apiV1Group, n, l)
	}
}
