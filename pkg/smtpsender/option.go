package smtpsender

type Option func(*MailSender)

func Host(host string) Option {
	return func(ms *MailSender) {
		ms.host = host
	}
}

func Port(port string) Option {
	return func(ms *MailSender) {
		ms.port = port
	}
}

func Username(u string) Option {
	return func(ms *MailSender) {
		ms.Username = u
	}
}

func Password(p string) Option {
	return func(ms *MailSender) {
		ms.Password = p
	}
}
