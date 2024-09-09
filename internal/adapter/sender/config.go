package sender

type Config struct {
	Host     string `env:"MAIL_SMTP_HOST" defaule:"localhost"`
	Port     int    `env:"MAIL_SMTP_PORT" default:"1025"`
	From     string `env:"MAIL_SENDER" default:"no-report@test.ru"`
	Password string `env:"MAIL_SENDER_PASSWORD"`
	Secure   bool   `env:"MAIL_SECURE" default:"0"`
}
