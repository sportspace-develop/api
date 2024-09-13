package rest

import (
	"sport-space/internal/adapter/models"
	"sport-space/pkg/email"
)

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

type tUpdTournamentRequest struct {
	Title string `json:"title"`
}

type tTournament struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type tGetTorunamentsResponse struct {
	Data []tTournament `json:"data"`
}

type tCreateTeam struct {
	Title string `json:"title"`
}

type tTeam struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type tGetTeamResponse struct {
	ID      uint      `json:"id"`
	Title   string    `json:"title"`
	Players []tPlayer `json:"players"`
}

type tUpdTeamRequest struct {
	Title   *string `json:"title"`
	Players *[]uint `json:"player_ids"`
}

type tUpdTeamResponse struct {
	ID      uint      `json:"id"`
	Title   string    `json:"title"`
	Players []tPlayer `json:"players"`
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

type tGetPlayersResponse struct {
	Data []tPlayer `json:"data"`
}

type tUpdatePlayer struct {
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	LastName   string `json:"lastname"`
}

type tApplication struct {
	ID              uint   `json:"id"`
	TournamentID    uint   `json:"tournament_id"`
	TournamentTitle string `json:"tournament_title"`
	Status          string `json:"status"`
}

type tNewApplicationRequest struct {
	TournamentID uint   `json:"tournament_id"`
	PlayerIDs    []uint `json:"player_ids"`
}

type tNewApplicationResponse struct {
	ID              uint      `json:"id"`
	TournamentID    uint      `json:"tournament_id"`
	TournamentTitle string    `json:"tournament_title"`
	Status          string    `json:"status"`
	Players         []tPlayer `json:"players"`
}

type applicationStatus string

var (
	submit applicationStatus = "submit"
	cancel applicationStatus = "cancel"
)

var applicationMapStatus = map[applicationStatus]models.ApplicationStatus{
	submit: models.InProgress,
	cancel: models.Canceled,
}

type tUpdApplicationStatusRequest struct {
	Status  *applicationStatus `json:"status" enums:"submit,cancel"`
	Players *[]uint            `json:"player_ids"`
}

type tUpdApplicationResponse struct {
	ID              uint      `json:"id"`
	TournamentID    uint      `json:"tournament_id"`
	TournamentTitle string    `json:"tournament_title"`
	Status          string    `json:"status"`
	Players         []tPlayer `json:"players"`
}

type tGetApplicationsTeamResponse struct {
	Data []tApplication `json:"data"`
}

type tGetApplicationResponse struct {
	ID              uint      `json:"id"`
	TournamentID    uint      `json:"tournament_id"`
	TournamentTitle string    `json:"tournament_title"`
	Status          string    `json:"status"`
	Players         []tPlayer `json:"players"`
}

type tTournamentApplication struct {
	ID        uint   `json:"id"`
	TeamID    uint   `json:"taem_id"`
	TeamTitle string `json:"team_title"`
	Status    string `json:"status"`
}

type tGetTournamentApplicationsResponse struct {
	Data []tTournamentApplication `json:"data"`
}

type tGetTorunamentApplicationResponse struct {
	ID        uint      `json:"id"`
	TeamID    uint      `json:"taem_id"`
	TeamTitle string    `json:"team_title"`
	Status    string    `json:"status"`
	Players   []tPlayer `json:"players"`
}

type applicationTournamentStatus string

var (
	accept applicationTournamentStatus = "accept"
	reject applicationTournamentStatus = "reject"
)

var applicationTournamentMapStatus = map[applicationTournamentStatus]models.ApplicationStatus{
	accept: models.Accepted,
	reject: models.Rejected,
}

type tUpdTournamentApplicationRequest struct {
	Status applicationTournamentStatus `json:"status" enums:"accept,reject"`
}

type tUpdTournamentApplicationResponse struct {
	ID        uint   `json:"id"`
	TeamID    uint   `json:"taem_id"`
	TeamTitle string `json:"team_title"`
	Status    string `json:"status"`
}
