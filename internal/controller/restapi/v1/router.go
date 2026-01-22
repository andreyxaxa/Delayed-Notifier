package v1

import (
	"github.com/andreyxaxa/Delayed-Notifier/internal/usecase"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewNotificationRoutes(apiV1Group fiber.Router, n usecase.Notification, l logger.Interface) {
	r := &V1{n: n, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	{
		// API
		apiV1Group.Post("/notify", r.createNotification)
		apiV1Group.Get("/notify/:id", r.getStatus)
		apiV1Group.Delete("/notify/:id", r.cancelNotification)
	}
}
