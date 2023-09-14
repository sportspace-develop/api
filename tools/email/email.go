package email

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sport-space-api/logger"
	"time"

	gomail "gopkg.in/mail.v2"
)

var (
	cfg      config
	log      *logger.Logger
	CMessage chan tMsg = make(chan tMsg, 2)
)

type Email string

func (e Email) String() string {
	return string(e)
}

type tMsg struct {
	To      string
	Subject string
	Body    string
}

type config interface {
	Sender() string
	Password() string
	Host() string
	Port() int
	Secure() bool
}

func sender() {
	for {
		m, ok := <-CMessage
		if !ok {
			log.ERROR("Failed read from msg chan")
			continue
		}
		SendEmail(m.To, m.Subject, m.Body)

	}
}

func Init(c config) {
	log = logger.New("email")
	log.INFO("init email")
	cfg = c
	go sender()
}

func (e *Email) IsValid() bool {

	return true
}

func AddMail(email, subject, body string) {
	m := tMsg{
		To:      email,
		Subject: subject,
		Body:    body,
	}

	marshaled, _ := json.MarshalIndent(m, "", "  ")
	log.INFO("Added to chan: ", string(marshaled))
	CMessage <- m
}

func SendEmail(to, subject, body string) (bool, error) {

	start := time.Now()
	m := gomail.NewMessage()

	host := cfg.Host()
	port := cfg.Port()
	sender := cfg.Sender()
	password := cfg.Password()
	secure := cfg.Secure()

	// Set E-Mail sender
	m.SetHeader("From", sender)

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", fmt.Sprintf("%s", body))

	// Settings for SMTP server
	d := gomail.NewDialer(host, port, sender, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: !secure}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		log.ERROR(err.Error())
		panic(err)
	}

	duration := time.Since(start).Seconds()
	marshaled, _ := json.MarshalIndent(map[string]interface{}{
		"duration": duration,
		"to":       to,
		"subject":  subject,
		"body":     body,
	}, "", "  ")
	log.INFO("Sended:", string(marshaled))
	return true, nil
}

func SendCodeToEmail(email string, code string) (bool, error) {
	start := time.Now()
	m := gomail.NewMessage()

	host := cfg.Host()
	port := cfg.Port()
	sender := cfg.Sender()
	password := cfg.Password()
	secure := cfg.Secure()

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
		log.ERROR(err.Error())
		panic(err)
	}

	duration := time.Since(start).Seconds()
	marshaled, _ := json.MarshalIndent(map[string]interface{}{
		"duration": duration,
		"to":       email,
		"subject":  "Auth code",
		"code":     code,
	}, "", "  ")
	log.INFO(string(marshaled))
	return true, nil
}
