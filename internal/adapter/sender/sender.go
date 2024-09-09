package sender

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

type Sender struct {
	log     *zap.Logger
	message chan tMsg
	cfg     Config
}

type option func(s *Sender)

func SetLogger(l *zap.Logger) option {
	return func(s *Sender) {
		s.log = l
	}
}

func New(cfg Config, options ...option) (*Sender, error) {
	s := &Sender{
		log:     zap.NewNop(),
		message: make(chan tMsg, 2),
		cfg:     cfg,
	}

	for _, opt := range options {
		opt(s)
	}

	go s.sender()

	return s, nil
}

type Email string

func (e Email) String() string {
	return string(e)
}

type tMsg struct {
	To      string
	Subject string
	Body    string
}

func (s *Sender) sender() {
	for {
		m, ok := <-s.message
		if !ok {
			s.log.Error("Failed read from msg chan")
			return
		}
		s.SendEmail(m.To, m.Subject, m.Body)
	}
}

func (s *Sender) AddMail(email, subject, body string) {
	m := tMsg{
		To:      email,
		Subject: subject,
		Body:    body,
	}

	s.log.Debug("Added to chan", zap.String("To", email), zap.String("Subject", subject), zap.String("Body", body))
	s.message <- m
}

func (s *Sender) SendEmail(to, subject, body string) (bool, error) {

	start := time.Now()
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", s.cfg.From)

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := gomail.NewDialer(s.cfg.Host, s.cfg.Port, s.cfg.From, s.cfg.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: !s.cfg.Secure}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		s.log.Error("failed send email", zap.Error(err))
		panic(err)
	}

	duration := time.Since(start).Seconds()
	marshaled, _ := json.MarshalIndent(map[string]interface{}{
		"duration": duration,
		"to":       to,
		"subject":  subject,
		"body":     body,
	}, "", "  ")
	s.log.Debug("Sended", zap.String("mail", string(marshaled)))
	return true, nil
}

func (s *Sender) SendCodeToEmail(email string, code string) (bool, error) {
	start := time.Now()
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", s.cfg.From)

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Auth code")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", fmt.Sprintf("Code %s", code))

	// Settings for SMTP server
	d := gomail.NewDialer(s.cfg.Host, s.cfg.Port, s.cfg.From, s.cfg.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: !s.cfg.Secure}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		s.log.Error("send email", zap.Error(err))
		return false, err
	}

	duration := time.Since(start).Seconds()
	s.log.Debug("sended otp", zap.Float64("duration", duration), zap.String("to", email), zap.String("code", code))
	return true, nil
}
