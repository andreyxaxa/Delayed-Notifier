package v1

import (
	"github.com/andreyxaxa/Delayed-Notifier/internal/usecase"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type V1 struct {
	n usecase.Notification
	l logger.Interface
	v *validator.Validate
}
