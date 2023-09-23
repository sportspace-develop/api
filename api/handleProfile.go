package api

import (
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type getProfileGameTypesResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"`
}

type getProfileResponse struct {
	Success       bool                           `json:"success"`
	GameTypes     []getProfileGameTypesResponse  `json:"game_types"`
	Organizations []getOrganizationDataResponse  `json:"organizations"`
	Tournaments   []getTournamentsDataResponse   `json:"tournaments"`
	Teams         []getTeamDataResponse          `json:"teams"`
	InviteStatus  map[string]string              `json:"invite_status"`
	InviteToTeam  []getPlayerInvitesDataResponse `json:"invites_to_team"`
	Player        getPlayerDataResponse          `json:"player"`
}

// @Summary profile
// @Schemes
// @Description get profile data
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} getProfileResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile [get]
func GetProfile(c *gin.Context) {
	session := sessions.New(c)

	userId := session.GetUserId()
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, responseError{
			Success: false,
			Error:   100,
			Message: GetMessageErr(100),
		})
		return
	}

	user, err := model.FindUserById(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	gameTypes, err := model.GetGameTypes()
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	gameTypesResult := []getProfileGameTypesResponse{}
	for _, gType := range gameTypes {
		gameTypesResult = append(gameTypesResult, getProfileGameTypesResponse{
			ID:    gType.ID,
			Title: gType.Title,
			Name:  gType.Name,
		})
	}

	organizations, err := model.GetOrganizationsByUserId(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	organizationsResult := []getOrganizationDataResponse{}

	for _, organization := range organizations {
		organizationsResult = append(organizationsResult, getOrganizationDataResponse{
			ID:      organization.ID,
			Title:   organization.Title,
			Address: organization.Address,
		})
	}

	teams, err := model.GetTeamsByUser(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	teamsResult := []getTeamDataResponse{}
	for _, team := range teams {
		teamsResult = append(teamsResult, getTeamDataResponse{
			ID:      team.ID,
			Title:   team.Title,
			DGameID: team.DGameID,
		})
	}

	invites, err := model.GetPlayerInviteToTeamByEmail(strings.ToLower(user.Email), model.TISended)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	var invitesToTeamResult []getPlayerInvitesDataResponse
	for _, invite := range invites {
		team, err := model.GetTeamById(invite.TeamID)
		if err != nil {
			responseErrorNumber(c, err, 500, http.StatusInternalServerError)
			return
		}
		invitesToTeamResult = append(invitesToTeamResult, getPlayerInvitesDataResponse{
			ID:        invite.ID,
			TeamID:    invite.TeamID,
			TeamTitle: team.Title,
			Status:    invite.Status.ToString(),
		})
	}

	player, err := model.GetPlayerFullByUserId(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	var playerBDay string
	if player.BDay.Valid {
		playerBDay = player.BDay.Time.Format(time.DateOnly)
	}

	var playerTeams []getPlayerTeamDataResponse
	for _, team := range player.Teams {
		playerTeams = append(playerTeams, getPlayerTeamDataResponse{
			ID:    team.ID,
			Title: team.Title,
		})
	}

	var tournamentsResult = []getTournamentsDataResponse{}
	tournamets, err := model.GetTournamentsByUser(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	for _, tournament := range tournamets {
		tournamentsResult = append(tournamentsResult, getTournamentsDataResponse{
			ID:                    tournament.ID,
			Title:                 tournament.Title,
			StartDate:             tournament.StartDate.Format(time.DateTime),
			EndDate:               tournament.EndDate.Format(time.DateTime),
			StartRegistrationDate: tournament.StartRegistrationDate.Format(time.DateTime),
			EndRegistrationDate:   tournament.EndRegistrationDate.Format(time.DateTime),
			IsTeam:                tournament.IsTeam,
			DGameID:               tournament.DGameID,
		})
	}

	playerResult := getPlayerDataResponse{
		FirstName:  player.FirstName,
		SecondName: player.SecondName,
		LastName:   player.LastName,
		BDay:       playerBDay,
		Teams:      playerTeams,
	}

	c.JSON(http.StatusOK, getProfileResponse{
		Success:       true,
		GameTypes:     gameTypesResult,
		Organizations: organizationsResult,
		Tournaments:   tournamentsResult,
		Teams:         teamsResult,
		InviteToTeam:  invitesToTeamResult,
		Player:        playerResult,
		InviteStatus:  model.InviteStatus,
	})
}
