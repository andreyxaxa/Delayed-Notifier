package workerpublisher

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/andreyxaxa/Delayed-Notifier/internal/usecase"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
)

const (
	_defaultTimeout    = 20 * time.Second
	_defaultTick       = 10 * time.Second
	_defaultMaxRetries = 3
)

type PublisherWorker struct {
	pub NotificationPublisher
	n   usecase.Notification
	l   logger.Interface

	ticker     *time.Ticker
	maxRetries int
	tick       time.Duration
	timeout    time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	started atomic.Bool
}

func New(publisher NotificationPublisher, n usecase.Notification, l logger.Interface, opts ...Option) *PublisherWorker {
	pub := &PublisherWorker{
		pub:        publisher,
		n:          n,
		l:          l,
		maxRetries: _defaultMaxRetries,
		tick:       _defaultTick,
		timeout:    _defaultTimeout,
	}

	for _, opt := range opts {
		opt(pub)
	}

	return pub
}

func (p *PublisherWorker) Start(ctx context.Context) error {
	if !p.started.CompareAndSwap(false, true) {
		return fmt.Errorf("workerpublisher - Start - worker already started")
	}

	p.ctx, p.cancel = context.WithCancel(ctx)
	p.ticker = time.NewTicker(p.tick)

	p.wg.Add(1)

	go func() {
		defer func() {
			p.ticker.Stop()
			p.wg.Done()
		}()

		for {
			select {
			case <-p.ctx.Done():
				return
			case t := <-p.ticker.C:
				batchCtx, batchCancel := context.WithTimeout(p.ctx, p.timeout)
				p.processBatch(batchCtx, t)
				batchCancel()
			}
		}
	}()

	return nil
}

func (p *PublisherWorker) processBatch(ctx context.Context, now time.Time) {
	ntfs, err := p.n.GetPendingNotifications(ctx, now)
	if err != nil {
		p.l.Error(err, "workerpublisher - PublisherWorker - processBatch")
		return
	}

	if len(ntfs) == 0 {
		return
	}

	for _, ntf := range ntfs {
		select {
		case <-ctx.Done():
			return
		default:
		}

		p.processNotification(ctx, ntf)
	}
}

func (p *PublisherWorker) processNotification(ctx context.Context, notification entity.Notification) {
	now := time.Now()

	if notification.RetryCount >= p.maxRetries {
		err := p.n.MarkAsFailed(ctx, notification.UID, now)
		if err != nil {
			p.l.Error(err, "workerpublisher - PublisherWorker - processNotification")
		}
		return
	}

	err := p.pub.PublishNotification(ctx, notification)
	if err != nil {
		p.l.Error(err, "workerpublisher - PublisherWorker - processNotification")

		nt := calculateNextTry(notification.RetryCount, now)
		err := p.n.RetryNotification(ctx, notification.UID, nt, now)
		if err != nil {
			p.l.Error(err, "workerpublisher - PublisherWorker - processNotification")
		}
		return
	}

	err = p.n.MarkAsProcessing(ctx, notification.UID, time.Now())
	if err != nil {
		p.l.Error(err, "workerpublisher - PublisherWorker - processNotification")

		return
	}
}

func (p *PublisherWorker) Shutdown(ctx context.Context) error {
	if !p.started.Load() {
		return nil
	}

	if p.cancel != nil {
		p.cancel()
	}

	done := make(chan struct{})

	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func calculateNextTry(retryCount int, now time.Time) time.Time {
	delay := time.Duration(1<<retryCount) * time.Minute

	return now.Add(delay)
}
