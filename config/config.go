package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	App AppConfig
)

type source string

const (
	DEV  source = "dev"
	TEST source = "test"
	PROD source = "prod"
)

type AppConfig struct {
	Source       source
	JWTSecret    string
	JWTLongTime  int
	CookieSecret string
}

func Init() {
	var err error
	godotenv.Load(".env")
	App = AppConfig{
		Source:       source(os.Getenv("SOURCE")),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		CookieSecret: os.Getenv("COOKIE_SECRET"),
	}
	App.JWTLongTime, err = strconv.Atoi(os.Getenv("JWT_LONG_TIME"))
	if err != nil {
		App.JWTLongTime = 600
	}
}

type DBCfg struct{}

func (cfg DBCfg) Host() string {
	return os.Getenv("DB_HOST")
}
func (cfg DBCfg) Port() string {
	return os.Getenv("DB_PORT")
}
func (cfg DBCfg) Username() string {
	return os.Getenv("DB_USERNAME")
}
func (cfg DBCfg) Password() string {
	return os.Getenv("DB_PASSWORD")
}
func (cfg DBCfg) DBName() string {
	return os.Getenv("DB_NAME")
}

type MailCfg struct{}

func (cfg MailCfg) Sender() string {
	return os.Getenv("MAIL_SENDER")
}
func (cfg MailCfg) Password() string {
	return os.Getenv("MAIL_SENDER_PASSWORD")
}
func (cfg MailCfg) Host() string {
	return os.Getenv("MAIL_SMTP_HOST")
}
func (cfg MailCfg) Port() int {
	port, err := strconv.Atoi(os.Getenv("MAIL_SMTP_PORT"))
	if err != nil {
		port = 1025
	}
	return port
}
func (cfg MailCfg) Secure() bool {
	secure, err := strconv.ParseBool(os.Getenv("MAIL_SECURE"))
	if err != nil {
		secure = false
	}
	return secure
}
