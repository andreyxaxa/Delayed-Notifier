package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	_defaultAddr        = "localhost:6379"
	_defaultDB          = 0
	_defaultMaxRetries  = 5
	_defaultUser        = ""
	_defaultPassword    = ""
	_defaultDialTimeout = 10 * time.Second
	_defaultTimeout     = 5 * time.Second
)

type Client struct {
	Client      *redis.Client
	addr        string
	password    string
	user        string
	db          int
	maxRetries  int
	dialTimeout time.Duration
	timeout     time.Duration
}

func New(ctx context.Context, opts ...Option) (*Client, error) {
	c := &Client{
		Client:      nil,
		addr:        _defaultAddr,
		user:        _defaultUser,
		password:    _defaultPassword,
		db:          _defaultDB,
		maxRetries:  _defaultMaxRetries,
		dialTimeout: _defaultDialTimeout,
		timeout:     _defaultTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	rc := redis.NewClient(&redis.Options{
		Addr:         c.addr,
		Username:     c.user,
		Password:     c.password,
		DB:           c.db,
		MaxRetries:   c.maxRetries,
		DialTimeout:  c.dialTimeout,
		ReadTimeout:  c.timeout,
		WriteTimeout: c.timeout,
	})

	if err := rc.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis - New - rc.Ping: %w", err)
	}

	c.Client = rc

	return c, nil
}
