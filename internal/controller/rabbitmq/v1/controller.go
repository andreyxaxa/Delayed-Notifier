package v1

import (
	"github.com/andreyxaxa/Delayed-Notifier/internal/usecase"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
)

type V1 struct {
	n usecase.Notification
	l logger.Interface
}
