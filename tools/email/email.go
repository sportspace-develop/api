package email

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

type Email string

func (e *Email) IsValid() bool {

	return true
}

func SendCodeToEmail(email string, code string) (bool, error) {

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "from@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", "to@example.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Auth code")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", fmt.Sprintf("Code %s", code))

	// Settings for SMTP server
	d := gomail.NewDialer("mail", 1025, "from@gmail.com", "<email_password>")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return true, nil
}
