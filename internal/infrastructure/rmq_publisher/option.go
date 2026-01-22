package rmqpublisher

type Option func(*NotificationPublisher)

func ContentType(contentType string) Option {
	return func(np *NotificationPublisher) {
		np.contentType = contentType
	}
}

func Exchange(exchange string) Option {
	return func(np *NotificationPublisher) {
		np.exchange = exchange
	}
}

func RoutingKey(rkey string) Option {
	return func(np *NotificationPublisher) {
		np.routingKey = rkey
	}
}
