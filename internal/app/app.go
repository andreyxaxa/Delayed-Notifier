package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyxaxa/Delayed-Notifier/config"
	"github.com/andreyxaxa/Delayed-Notifier/internal/controller/restapi"
	mailsender "github.com/andreyxaxa/Delayed-Notifier/internal/infrastructure/mail_sender"
	telegramsender "github.com/andreyxaxa/Delayed-Notifier/internal/infrastructure/telegram_sender"
	"github.com/andreyxaxa/Delayed-Notifier/internal/repo/cache"
	"github.com/andreyxaxa/Delayed-Notifier/internal/repo/persistent"
	"github.com/andreyxaxa/Delayed-Notifier/internal/usecase/notification"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/httpserver"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/postgres"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/redis"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/smtpsender"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/telegrambotsender"
)

func Run(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %v", err))
	}
	defer pg.Close()

	// Cache
	rd, err := redis.New(ctx, redis.Addr(cfg.Redis.Addr), redis.DB(cfg.Redis.DB))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis.New: %v", err))
	}

	// SMTP Sender
	smtpSender := smtpsender.New(
		smtpsender.Host(cfg.SMTPMail.Host),
		smtpsender.Port(cfg.SMTPMail.Port),
		smtpsender.Username(cfg.SMTPMail.Username),
		smtpsender.Password(cfg.SMTPMail.Password),
	)

	// Telegram bot-api Sender
	tgSender, err := telegrambotsender.New(cfg.Telegram.Token)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - telegrambotsender.New: %v", err))
	}

	// Use-Case
	notificationUseCase := notification.New(
		persistent.New(pg),
		mailsender.New(smtpSender),
		telegramsender.New(tgSender),
		cache.New(rd),
		l,
	)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpServer.App, notificationUseCase, l)

	// Start servers
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %v", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %v", err))
	}
}
