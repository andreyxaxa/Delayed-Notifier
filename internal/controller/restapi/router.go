package restapi

import (
	"github.com/andreyxaxa/Delayed-Notifier/config"
	_ "github.com/andreyxaxa/Delayed-Notifier/docs" // Swagger docs.
	v1 "github.com/andreyxaxa/Delayed-Notifier/internal/controller/restapi/v1"
	"github.com/andreyxaxa/Delayed-Notifier/internal/usecase"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title 	 Delayed Notifier
// @version  1.0
// @host 	 localhost:8080
// @BasePath /v1
func NewRouter(app *fiber.App, cfg *config.Config, n usecase.Notification, l logger.Interface) {
	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewNotificationRoutes(apiV1Group, n, l)
	}
}
