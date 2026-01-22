package rmqserver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/pkg/logger"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

type CallHandler func(context.Context, amqp091.Delivery) error

var (
	ErrNoHandler = errors.New("no handler")
)

const (
	_defaultConnectTimeout = 5 * time.Second
	_defaultHeartbeat      = 10 * time.Second
	_defaultConnectionName = "myconnection"
	_defaultQueueName      = "v1.sendNotification"
	_defaultExchangeName   = "notifications"
	_defaultRoutingKey     = "v1.sendNotification"
	_defaultConsumerTag    = "delayed-notifier"
	_defaultWorkers        = 1
	_defaultPrefetchCount  = 1

	_defaultAttempts = 3
	_defaultDelay    = 3 * time.Second
	_defaultBackoff  = 2
)

type Server struct {
	Client      *rabbitmq.RabbitClient
	router      map[string]CallHandler
	queue       string
	exchange    string
	routingKey  string
	consumerTag string

	prefetch int
	workers  int
	notify   chan error

	l logger.Interface
}

func New(url string, router map[string]CallHandler, l logger.Interface, opts ...Option) (*Server, error) {
	strat := retry.Strategy{
		Attempts: _defaultAttempts,
		Delay:    _defaultDelay,
		Backoff:  _defaultBackoff,
	}

	cfg := rabbitmq.ClientConfig{
		URL:            url,
		ConnectionName: _defaultConnectionName,
		ConnectTimeout: _defaultConnectTimeout,
		Heartbeat:      _defaultHeartbeat,
		ReconnectStrat: strat,
		ConsumingStrat: strat,
		ProducingStrat: strat,
	}

	client, err := rabbitmq.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq - server - New - rabbitmq.NewClient: %w", err)
	}

	s := &Server{
		Client:      client,
		router:      router,
		queue:       _defaultQueueName,
		exchange:    _defaultExchangeName,
		routingKey:  _defaultRoutingKey,
		consumerTag: _defaultConsumerTag,
		prefetch:    _defaultPrefetchCount,
		workers:     _defaultWorkers,
		notify:      make(chan error, 1),
		l:           l,
	}

	for _, opt := range opts {
		opt(s)
	}

	err = s.declareQueue()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq - server - New - s.declareQueue: %w", err)
	}

	return s, nil
}

func (s *Server) declareQueue() error {
	err := s.Client.DeclareExchange(
		s.exchange,
		"direct",
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rabbitmq - server - declare - s.client.DeclareExchange: %w", err)
	}

	err = s.Client.DeclareQueue(
		s.queue,
		s.exchange,
		s.routingKey,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rabbitmq - server - declare - s.client.DeclareQueue: %w", err)
	}

	return nil
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		cfg := rabbitmq.ConsumerConfig{
			Queue:         s.queue,
			ConsumerTag:   s.consumerTag,
			AutoAck:       false,
			Workers:       s.workers,
			PrefetchCount: s.prefetch,
			Ask: rabbitmq.AskConfig{
				Multiple: false,
			},
			Nack: rabbitmq.NackConfig{
				Multiple: false,
				Requeue:  true,
			},
		}

		handler := func(ctx context.Context, d amqp091.Delivery) error {
			callHandler, ok := s.router[d.RoutingKey]
			if !ok {
				s.l.Error(ErrNoHandler, "rabbitmq - server - Start")

				return ErrNoHandler
			}

			err := callHandler(ctx, d)
			if err != nil {
				s.l.Error(err, "rabbitmq - server - Start - callHandler")

				return err
			}
			return nil
		}

		consumer := rabbitmq.NewConsumer(s.Client, cfg, handler)

		err := consumer.Start(ctx)
		if err != nil {
			s.notify <- err
			close(s.notify)
		}
	}()

	s.l.Info("rabbitmq server - Server - Started")
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	err := s.Client.Close()
	if err != nil {
		return err
	}

	s.l.Info("rabbitmq server - Server - Shutdown")

	return nil
}
