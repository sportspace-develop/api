package rest

import (
	"time"

	"sport-space/internal/adapter/models"
	"sport-space/pkg/email"
)

type sportTime string

func (st *sportTime) DateTime() *time.Time {
	if st == nil {
		return nil
	}

	res, _ := time.Parse(defaultDateTimeFormat, string(*st))
	return &res
}

func (st *sportTime) Date() *time.Time {
	if st == nil {
		return nil
	}

	res, _ := time.Parse(defaultDateFormat, string(*st))
	return &res
}

func formatDateTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	_t := t.Format(defaultDateTimeFormat)
	return _t
}

func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	_t := t.Format(defaultDateFormat)
	return _t
}

func (st *sportTime) String() string {
	if st == nil {
		return ""
	}
	return string(*st)
}

type pagination struct {
	TotalRecords int  `json:"total_records"`
	CurrentPage  int  `json:"current_page"`
	TotalPages   int  `json:"total_pages"`
	NextPage     *int `json:"next_page"`
	PrevPage     *int `json:"prev_page"`
	Limit        int  `json:"-"`
	StartRow     int  `json:"-"`
	EndRow       int  `json:"-"`
}

type tAuthorization struct {
	Email string `json:"email"`
	// Password string `json:"password"`
	OTP string `json:"otp"`
}

type tRequestOTP struct {
	Email email.Email `json:"email"`
}

type tCreateTournamentRequest struct {
	Title             string     `json:"title" validate:"required"`
	StartDate         *sportTime `json:"start_date" example:"2024-12-31 00:00:00" validate:"required"`
	EndDate           *sportTime `json:"end_date" example:"2024-12-31 00:00:00" validate:"required"`
	RegisterStartDate *sportTime `json:"register_start_date" example:"2024-12-31 00:00:00"`
	RegisterEndDate   *sportTime `json:"register_end_date" example:"2024-12-31 00:00:00"`
	LogoURL           string     `json:"logo_url"`
}

func (tct tCreateTournamentRequest) IsValid() bool {
	return !(tct.Title == "" || tct.StartDate == nil || tct.EndDate == nil)
}

type tUpdTournamentRequest struct {
	Title             string     `json:"title"`
	StartDate         *sportTime `json:"start_date" example:"2024-12-31 00:00:00" validate:"required"`
	EndDate           *sportTime `json:"end_date" example:"2024-12-31 00:00:00" validate:"required"`
	RegisterStartDate *sportTime `json:"register_start_date" example:"2024-12-31 00:00:00"`
	RegisterEndDate   *sportTime `json:"register_end_date" example:"2024-12-31 00:00:00"`
	LogoURL           string     `json:"logo_url"`
}

func (tutr tUpdTournamentRequest) IsValid() bool {
	return !(tutr.Title == "" || tutr.StartDate == nil || tutr.EndDate == nil)
}

type tTournamentResponse struct {
	ID                uint   `json:"id"`
	Title             string `json:"title"`
	StartDate         string `json:"start_date" example:"2024-12-31 00:00:00"`
	EndDate           string `json:"end_date" example:"2024-12-31 00:00:00"`
	RegisterStartDate string `json:"register_start_date" example:"2024-12-31 00:00:00"`
	RegisterEndDate   string `json:"register_end_date" example:"2024-12-31 00:00:00"`
	LogoURL           string `json:"logo_url"`
	LogoExternalURL   string `json:"logo_external_url"`
}

type tGetTorunamentsResponse struct {
	Pagination pagination            `json:"pagination"`
	Data       []tTournamentResponse `json:"data"`
}

type tCreateTeam struct {
	Title    string `json:"title"`
	LogotURL string `json:"logo_url"`
	PhotoURL string `json:"photo_url"`
}

type tTeam struct {
	ID               uint   `json:"id"`
	Title            string `json:"title"`
	LogoURL          string `json:"logo_url"`
	LogoExternalURL  string `json:"logo_external_url"`
	PhotoURL         string `json:"photo_url"`
	PhotoExternalURL string `json:"photo_external_url"`
	CreatedAt        string `json:"created_at"`
}

type tGetTeamsResponse struct {
	Pagination pagination `json:"pagination"`
	Data       []tTeam    `json:"data"`
}

type tGetTeamResponse struct {
	ID               uint              `json:"id"`
	Title            string            `json:"title"`
	Players          []tPlayerResponse `json:"players"`
	LogoURL          string            `json:"logo_url"`
	LogoExternalURL  string            `json:"logo_external_url"`
	PhotoURL         string            `json:"photo_url"`
	PhotoExternalURL string            `json:"photo_external_url"`
	CreatedAt        string            `json:"created_at"`
}

type tUpdTeamRequest struct {
	Title    string  `json:"title"`
	LogotURL string  `json:"logo_url"`
	PhotoURL string  `json:"photo_url"`
	Players  *[]uint `json:"player_ids"`
}

type tUpdTeamResponse struct {
	ID               uint               `json:"id"`
	Title            string             `json:"title"`
	Players          *[]tPlayerResponse `json:"players"`
	LogoURL          string             `json:"logo_url"`
	LogoExternalURL  string             `json:"logo_external_url"`
	PhotoURL         string             `json:"photo_url"`
	PhotoExternalURL string             `json:"photo_external_url"`
	CreatedAt        string             `json:"created_at"`
}

type tNewPlayerRequest struct {
	FirstName  string     `json:"firstname"`
	SecondName string     `json:"secondname"`
	LastName   string     `json:"lastname"`
	PhotoURL   string     `json:"photo_url"`
	BDay       *sportTime `json:"b_day" example:"2024-12-31"`
}

func (tnp tNewPlayerRequest) IsValid() bool {
	return !(tnp.FirstName == "" || tnp.LastName == "")
}

type tPlayerResponse struct {
	ID               uint   `json:"id"`
	FirstName        string `json:"firstname"`
	SecondName       string `json:"secondname"`
	LastName         string `json:"lastname"`
	PhotoURL         string `json:"photo_url"`
	PhotoExternalURL string `json:"photo_external_url"`
	BDay             string `json:"b_day" example:"2024-12-31"`
}

type tGetPlayersResponse struct {
	Pagination pagination        `json:"pagination"`
	Data       []tPlayerResponse `json:"data"`
}

type tUpdatePlayerRequest struct {
	FirstName  string     `json:"firstname"`
	SecondName string     `json:"secondname"`
	LastName   string     `json:"lastname"`
	PhotoURL   string     `json:"photo_url"`
	BDay       *sportTime `json:"b_day" example:"2024-12-31"`
}

func (tup tUpdatePlayerRequest) IsValid() bool {
	return !(tup.FirstName == "" || tup.LastName == "")
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
	ID              uint              `json:"id"`
	TournamentID    uint              `json:"tournament_id"`
	TournamentTitle string            `json:"tournament_title"`
	Status          string            `json:"status"`
	Players         []tPlayerResponse `json:"players"`
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
	ID              uint              `json:"id"`
	TournamentID    uint              `json:"tournament_id"`
	TournamentTitle string            `json:"tournament_title"`
	Status          string            `json:"status"`
	Players         []tPlayerResponse `json:"players"`
}

type tGetApplicationsTeamResponse struct {
	Data []tApplication `json:"data"`
}

type tGetApplicationResponse struct {
	ID              uint              `json:"id"`
	TournamentID    uint              `json:"tournament_id"`
	TournamentTitle string            `json:"tournament_title"`
	Status          string            `json:"status"`
	Players         []tPlayerResponse `json:"players"`
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
	ID        uint              `json:"id"`
	TeamID    uint              `json:"taem_id"`
	TeamTitle string            `json:"team_title"`
	Status    string            `json:"status"`
	Players   []tPlayerResponse `json:"players"`
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

type tHandlerUploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}
