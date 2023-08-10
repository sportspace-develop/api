package email

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

var (
	cfg config
)

type config interface {
	Sender() string
	Password() string
	Host() string
	Port() int
	Secure() bool
}

func Init(c config) {
	cfg = c
}

type Email string

func (e *Email) IsValid() bool {

	return true
}

func SendCodeToEmail(email string, code string) (bool, error) {

	m := gomail.NewMessage()

	host := cfg.Host()
	port := cfg.Port()
	sender := cfg.Sender()
	password := cfg.Password()
	secure := cfg.Secure()

	fmt.Println(host, port, sender, password, secure)

	// Set E-Mail sender
	m.SetHeader("From", sender)

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Auth code")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", fmt.Sprintf("Code %s", code))

	// Settings for SMTP server
	d := gomail.NewDialer(host, port, sender, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: !secure}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return true, nil
}
