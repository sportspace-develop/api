package email

type Email string

func (e *Email) IsValid() bool {

	return true
}

func SendCodeToEmail(email string, code string) (bool, error) {

	return true, nil
}
