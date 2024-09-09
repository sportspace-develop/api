package email

import "net/mail"

type Email string

func (e Email) String() string {
	return string(e)
}

func (e Email) IsValid() bool {
	_, err := mail.ParseAddress(e.String())
	return err == nil
}
