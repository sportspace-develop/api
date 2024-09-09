package rest

import "sport-space/pkg/email"

type tAuthorization struct {
	Email string `json:"email"`
	// Password string `json:"password"`
	OTP string `json:"otp"`
}

type tRequestOTP struct {
	Email email.Email `json:"email"`
}

type tCreateTournament struct {
	Title string `json:"title"`
}

type tUpdateTournament struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type tTournament struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type tCreateTeam struct {
	Title string `json:"title"`
}

type tTeam struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type tUpdateTeam struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type tNewPlayer struct {
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	LastName   string `json:"lastname"`
}

type tPlayer struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	LastName   string `json:"lastname"`
}

type tUpdatePlayer struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	LastName   string `json:"lastname"`
}

type tNewApplication struct {
	TournamentID uint `json:"tournament_id"`
}
