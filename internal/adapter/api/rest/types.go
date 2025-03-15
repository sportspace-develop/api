package rest

import (
	"time"

	"sport-space/internal/adapter/models"
	"sport-space/pkg/email"
)

var (
	defaultDateTimeFormat = time.RFC3339 // "2006-01-02T15:04:05+03:00"
	defaultDateFormat     = time.RFC3339 //"2006-01-02T15:04:05+05:00"
)

type sportTime string

func (st *sportTime) DateTime() *time.Time {
	if st == nil || *st == "" {
		return nil
	}

	res, _ := time.Parse(defaultDateTimeFormat, string(*st))
	return &res
}

func (st *sportTime) Date() *time.Time {
	if st == nil || *st == "" {
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
	TotalRecords int  `json:"totalRecords"`
	CurrentPage  int  `json:"currentPage"`
	TotalPages   int  `json:"totalPages"`
	NextPage     *int `json:"nextPage"`
	PrevPage     *int `json:"prevPage"`
	Limit        int  `json:"-"`
	StartRow     int  `json:"-"`
	EndRow       int  `json:"-"`
}

type tLoginResponse struct {
	UserID uint   `json:"userID" example:"1"`
	Error  string `json:"error"`
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
	Description       string     `json:"description"`
	Organization      string     `json:"organization"`
	StartDate         *sportTime `json:"startDate" example:"2024-12-31T06:00:00+03:00" validate:"required"`
	EndDate           *sportTime `json:"endDate" example:"2024-12-31T06:00:00+03:00" validate:"required"`
	RegisterStartDate *sportTime `json:"registerStartDate" example:"2024-12-31T06:00:00+03:00"`
	RegisterEndDate   *sportTime `json:"registerEndDate" example:"2024-12-31T06:00:00+03:00"`
	LogoURL           string     `json:"logoUrl"`
}

func (tct tCreateTournamentRequest) IsValid() bool {
	return !(tct.Title == "" || tct.StartDate == nil || tct.EndDate == nil)
}

type tUpdTournamentRequest struct {
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	Organization      string     `json:"organization"`
	StartDate         *sportTime `json:"startDate" example:"2024-12-31T06:00:00+03:00" validate:"required"`
	EndDate           *sportTime `json:"endDate" example:"2024-12-31T06:00:00+03:00" validate:"required"`
	RegisterStartDate *sportTime `json:"registerStartDate" example:"2024-12-31T06:00:00+03:00"`
	RegisterEndDate   *sportTime `json:"registerEndDate" example:"2024-12-31T06:00:00+03:00"`
	LogoURL           string     `json:"logoUrl"`
}

func (tutr tUpdTournamentRequest) IsValid() bool {
	return !(tutr.Title == "" || tutr.StartDate == nil || tutr.EndDate == nil)
}

type tTournamentResponse struct {
	ID                uint   `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Organization      string `json:"organization"`
	OrganizationID    uint   `json:"organizationID"`
	StartDate         string `json:"startDate" example:"2024-12-31T06:00:00+03:00"`
	EndDate           string `json:"endDate" example:"2024-12-31T06:00:00+03:00"`
	RegisterStartDate string `json:"registerStartDate" example:"2024-12-31T06:00:00+03:00"`
	RegisterEndDate   string `json:"registerEndDate" example:"2024-12-31T06:00:00+03:00"`
	LogoURL           string `json:"logoUrl"`
}

type tGetTorunamentsResponse struct {
	Pagination pagination            `json:"pagination"`
	Data       []tTournamentResponse `json:"data"`
}

type tCreateTeam struct {
	Title    string `json:"title"`
	LogoURL  string `json:"logoUrl"`
	PhotoURL string `json:"photoUrl"`
}

type tTeam struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	LogoURL   string `json:"logoUrl"`
	PhotoURL  string `json:"photoUrl"`
	CreatedAt string `json:"createdAt"`
}

type tGetTeamsResponse struct {
	Pagination pagination `json:"pagination"`
	Data       []tTeam    `json:"data"`
}

type tGetTeamResponse struct {
	ID        uint              `json:"id"`
	Title     string            `json:"title"`
	Players   []tPlayerResponse `json:"players"`
	LogoURL   string            `json:"logoUrl"`
	PhotoURL  string            `json:"photoUrl"`
	CreatedAt string            `json:"createdAt"`
}

type tUpdTeamRequest struct {
	Title    string  `json:"title"`
	LogoURL  string  `json:"logoUrl"`
	PhotoURL string  `json:"photoUrl"`
	Players  *[]uint `json:"playerIds"`
}

type tUpdTeamResponse struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Players   *[]tPlayerResponse `json:"players"`
	LogoURL   string             `json:"logoUrl"`
	PhotoURL  string             `json:"photoUrl"`
	CreatedAt string             `json:"createdAt"`
}

type tNewPlayerRequest struct {
	FirstName  string     `json:"firstName"`
	SecondName string     `json:"secondName"`
	LastName   string     `json:"lastName"`
	PhotoURL   string     `json:"photoUrl"`
	BDay       *sportTime `json:"bDay" example:"2024-12-31T06:00:00+03:00"`
}

func (tnp tNewPlayerRequest) IsValid() bool {
	return !(tnp.FirstName == "" || tnp.LastName == "")
}

type tPlayerResponse struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	LastName   string `json:"lastName"`
	PhotoURL   string `json:"photoUrl"`
	BDay       string `json:"bDay" example:"2024-12-31T06:00:00+03:00"`
}

type tNewPlayerBatchRequest struct {
	FirstName  string     `json:"firstName"`
	SecondName string     `json:"secondName"`
	LastName   string     `json:"lastName"`
	PhotoURL   string     `json:"photoUrl"`
	BDay       *sportTime `json:"bDay" example:"2024-12-31T06:00:00+03:00"`
	ID         uint       `json:"id"`
}

func (tnp tNewPlayerBatchRequest) IsValid() bool {
	return !(tnp.FirstName == "" || tnp.LastName == "")
}

type tPlayerBatchResponse struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	LastName   string `json:"lastName"`
	PhotoURL   string `json:"photoUrl"`
	BDay       string `json:"bDay" example:"2024-12-31T06:00:00+03:00"`
}

type tNewPlayerBatchResponse struct {
	// Pagination pagination         `json:"pagination"`
	Data []tPlayerBatchResponse `json:"data"`
}

type tGetPlayersResponse struct {
	Pagination pagination        `json:"pagination"`
	Data       []tPlayerResponse `json:"data"`
}

type tUpdatePlayerRequest struct {
	FirstName  string     `json:"firstName"`
	SecondName string     `json:"secondName"`
	LastName   string     `json:"lastName"`
	PhotoURL   string     `json:"photoUrl"`
	BDay       *sportTime `json:"bDay" example:"2024-12-31T06:00:00+03:00"`
}

func (tup tUpdatePlayerRequest) IsValid() bool {
	return !(tup.FirstName == "" || tup.LastName == "")
}

type tApplication struct {
	ID                uint   `json:"id"`
	TournamentID      uint   `json:"tournamentId"`
	TournamentTitle   string `json:"tournamentTitle"`
	TournamentLogoURL string `json:"tournamentLogoUrl"`
	Status            string `json:"status"`
}

type tNewApplicationRequest struct {
	TournamentID uint   `json:"tournamentId"`
	PlayerIDs    []uint `json:"playerIds"`
}

type tNewApplicationResponse struct {
	ID              uint              `json:"id"`
	TournamentID    uint              `json:"tournamentId"`
	TournamentTitle string            `json:"tournamentTitle"`
	Status          string            `json:"status"`
	Players         []tPlayerResponse `json:"players"`
}

type applicationStatus string

var (
	submit applicationStatus = "submit"
	cancel applicationStatus = "cancel"
	draft  applicationStatus = "draft"
)

var applicationMapStatus = map[applicationStatus]models.ApplicationStatus{
	submit: models.InProgress,
	cancel: models.Canceled,
	draft:  models.Draft,
}

type tUpdApplicationStatusRequest struct {
	Status  *applicationStatus `json:"status" enums:"submit,cancel,draft"`
	Players *[]uint            `json:"playerIds"`
}

type tUpdApplicationResponse struct {
	ID              uint              `json:"id"`
	TournamentID    uint              `json:"tournamentId"`
	TournamentTitle string            `json:"tournamentTitle"`
	Status          string            `json:"status"`
	Players         []tPlayerResponse `json:"players"`
}

type tGetApplicationsTeamResponse struct {
	Data []tApplication `json:"data"`
}

type tGetApplicationResponse struct {
	ID              uint              `json:"id"`
	TournamentID    uint              `json:"tournamentId"`
	TournamentTitle string            `json:"tournamentTitle"`
	Status          string            `json:"status"`
	Players         []tPlayerResponse `json:"players"`
}

type tTournamentApplication struct {
	ID          uint   `json:"id"`
	TeamID      uint   `json:"teamId"`
	TeamTitle   string `json:"teamTitle"`
	TeamLogoURL string `json:"teamLogoUrl"`
	Status      string `json:"status"`
}

type tGetTournamentApplicationsResponse struct {
	Data []tTournamentApplication `json:"data"`
}

type tGetTorunamentApplicationResponse struct {
	ID        uint              `json:"id"`
	TeamID    uint              `json:"teamId"`
	TeamTitle string            `json:"teamTitle"`
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
	TeamID    uint   `json:"teamId"`
	TeamTitle string `json:"teamTitle"`
	Status    string `json:"status"`
}

type tHandlerUploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}
