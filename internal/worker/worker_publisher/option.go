package workerpublisher

import "time"

type Option func(*PublisherWorker)

func Timeout(timeout time.Duration) Option {
	return func(pw *PublisherWorker) {
		pw.timeout = timeout
	}
}

func Tick(tick time.Duration) Option {
	return func(pw *PublisherWorker) {
		pw.tick = tick
	}
}

func MaxRetries(retries int) Option {
	return func(pw *PublisherWorker) {
		pw.maxRetries = retries
	}
}
