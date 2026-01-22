package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		HTTP     HTTP
		PG       PG
		Redis    Redis
		Log      Log
		RMQ      RMQ
		Worker   Worker
		SMTPMail SMTPMail
		Telegram Telegram
	}

	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	Redis struct {
		Addr        string `env:"REDIS_ADDR,required"`
		DB          int    `env:"REDIS_DB,required"`
		User        string `env:"REDIS_USER"`
		Password    string `env:"REDIS_PASSWORD"`
		DialTimeout int    `env:"REDIS_DIAL_TIMEOUT"`
		Timeout     int    `env:"REDIS_TIMEOUT"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	RMQ struct {
		URL           string `env:"RMQ_URL,required"`
		Queue         string `env:"RMQ_QUEUE"`
		Exchange      string `env:"RMQ_EXCHANGE"`
		RoutingKey    string `env:"RMQ_ROUTING_KEY"`
		Workers       int    `env:"RMQ_WORKERS"`
		PrefetchCount int    `env:"RMQ_PREFETCH_COUNT"`
	}

	Worker struct {
		ShutdownTimeout int `env:"WORKER_SHUTDOWN_TIMEOUT,required"`
	}

	SMTPMail struct {
		Username string `env:"SMTPMAIL_USERNAME,required"`
		Password string `env:"SMTPMAIL_PASSWORD,required"`
		Host     string `env:"SMTPMAIL_HOST,required"`
		Port     string `env:"SMTPMAIL_PORT,required"`
	}

	Telegram struct {
		Token string `env:"TELEGRAM_BOT_TOKEN,required"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %v", err)
	}

	return cfg, nil
}
