package rmqserver

type Option func(s *Server)

func Queue(queue string) Option {
	return func(s *Server) {
		s.queue = queue
	}
}

func Exchange(exchange string) Option {
	return func(s *Server) {
		s.exchange = exchange
	}
}

func RoutingKey(key string) Option {
	return func(s *Server) {
		s.routingKey = key
	}
}

func Workers(workers int) Option {
	return func(s *Server) {
		s.workers = workers
	}
}

func PrefetchCount(count int) Option {
	return func(s *Server) {
		s.prefetch = count
	}
}

// TODO: добавить в конфиг
func ConsumerTag(tag string) Option {
	return func(s *Server) {
		s.consumerTag = tag
	}
}
