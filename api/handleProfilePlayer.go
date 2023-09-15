package api

import (
	"database/sql"
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type getPlayerTeamDataResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type getPlayerDataResponse struct {
	FirstName  string                      `json:"first_name"`
	SecondName string                      `json:"second_name"`
	LastName   string                      `json:"last_name"`
	BDay       string                      `json:"b_day" example:"2000-12-31"`
	Teams      []getPlayerTeamDataResponse `json:"teams"`
}

type getPlayerResponse struct {
	Success bool                  `json:"success"`
	Data    getPlayerDataResponse `json:"data"`
}

// @Summary profile player
// @Schemes
// @Description profile player
// @Tags profile player
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param team_id query int false "Team id"
// @Success 200 {object} getPlayerResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/player [get]
func GetPlayer(c *gin.Context) {
	session := sessions.New(c)

	userId := session.GetUserId()
	if userId == 0 {
		responseErrorNumber(c, nil, 100, http.StatusUnauthorized)
		return
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

	result := getPlayerDataResponse{
		FirstName:  player.FirstName,
		SecondName: player.SecondName,
		LastName:   player.LastName,
		BDay:       playerBDay,
	}

	for _, team := range player.Teams {
		result.Teams = append(result.Teams, getPlayerTeamDataResponse{
			ID:    team.ID,
			Title: team.Title,
		})
	}

	c.JSON(http.StatusOK, getPlayerResponse{
		Success: true,
		Data:    result,
	})
}

type updatePlayerRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
	BDay       string `json:"bday"`
}

type updatePlayerDataResponse struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
	BDay       string `json:"bday"`
}

type updatePlayerResponse struct {
	Success bool                     `json:"success"`
	Data    updatePlayerDataResponse `json:"data"`
}

// @Summary update profile player
// @Schemes
// @Description update profile player
// @Tags profile player
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body updatePlayerRequest false "Body"
// @Success 200 {object} updatePlayerResponse
// @Failure 400 {object} responseError
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/player [put]
func UpdatePlayer(c *gin.Context) {
	session := sessions.New(c)

	userId := session.GetUserId()
	if userId == 0 {
		responseErrorNumber(c, nil, 100, http.StatusUnauthorized)
		return
	}

	jsonData := updatePlayerRequest{}
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		responseErrorNumber(c, err, 1400, http.StatusInternalServerError)
		return
	}

	player, err := model.GetPlayerByUserId(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	if jsonData.BDay != "" {
		jsonDataBDay, err := time.Parse(time.DateOnly, jsonData.BDay)
		if err != nil {
			responseErrorNumber(c, err, 1401, http.StatusInternalServerError)
			return
		}
		var playerBDaySQL sql.NullTime
		err = playerBDaySQL.Scan(jsonDataBDay)
		if err != nil {
			responseErrorNumber(c, err, 1402, http.StatusInternalServerError)
			return
		}
		player.BDay = playerBDaySQL
	}

	player.FirstName = jsonData.FirstName
	player.SecondName = jsonData.SecondName
	player.LastName = jsonData.LastName
	player, err = model.UpdatePlayer(player)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	result := updatePlayerDataResponse{
		FirstName:  player.FirstName,
		SecondName: player.SecondName,
		LastName:   player.LastName,
		BDay:       player.BDay.Time.Format(time.DateOnly),
	}

	c.JSON(http.StatusOK, updatePlayerResponse{
		Success: true,
		Data:    result,
	})
}

type getPlayerInvitesDataResponse struct {
	ID        uint   `json:"id"`
	TeamID    uint   `json:"team_id"`
	TeamTitle string `json:"team_title"`
	Status    string `json:"status"`
}

type getPlayerInvitesResponse struct {
	Success      bool                           `json:"success"`
	InviteStatus map[string]string              `json:"invite_status"`
	Data         []getPlayerInvitesDataResponse `json:"data"`
}

// @Summary profile player invites
// @Schemes
// @Description profile player invites
// @Tags profile player
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param invite_id query int false "Invite id"
// @Success 200 {object} getPlayerInvitesResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/player/invite [get]
func GetPlayerInvite(c *gin.Context) {
	session := sessions.New(c)

	userId := session.GetUserId()
	if userId == 0 {
		responseErrorNumber(c, nil, 100, http.StatusUnauthorized)
		return
	}

	user, err := model.FindUserById(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	var inviteId int
	inviteId, _ = strconv.Atoi(c.DefaultQuery("invite_id", "0"))

	var invites []model.TeamInvite
	if inviteId > 0 {
		invite, err := model.GetPlayerInviteToTeamById(uint(inviteId), model.TISended)
		if err != nil {
			responseErrorNumber(c, err, 500, http.StatusInternalServerError)
			return
		}
		if invite.ID > 0 {
			invites = append(invites, invite)
		}
	} else {
		invites, err = model.GetPlayerInviteToTeamByEmail(strings.ToLower(user.Email), model.TISended)
		if err != nil {
			responseErrorNumber(c, err, 500, http.StatusInternalServerError)
			return
		}
	}

	var result []getPlayerInvitesDataResponse
	for _, invite := range invites {
		team, err := model.GetTeamById(invite.TeamID)
		if err != nil {
			responseErrorNumber(c, err, 500, http.StatusInternalServerError)
			return
		}
		result = append(result, getPlayerInvitesDataResponse{
			ID:        invite.ID,
			TeamID:    invite.TeamID,
			TeamTitle: team.Title,
			Status:    invite.Status.ToString(),
		})
	}

	c.JSON(http.StatusOK, getPlayerInvitesResponse{
		Success:      true,
		InviteStatus: model.InviteStatus,
		Data:         result,
	})
}

type updatePlayerInviteRequest struct {
	InviteID uint   `json:"invite_id"`
	Status   string `json:"status"`
	// Code     string `json:"code"`
}

// @Summary update profile player invites
// @Schemes
// @Description update profile player invites
// @Tags profile player
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body updatePlayerInviteRequest true "json"
// @Success 200 {object} responseSuccess
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/player/invite [put]
func UpdatePlayerInvite(c *gin.Context) {
	var err error
	session := sessions.New(c)

	userId := session.GetUserId()
	if userId == 0 {
		responseErrorNumber(c, nil, 100, http.StatusUnauthorized)
		return
	}

	user, err := model.FindUserById(userId)
	if err != nil {
		responseErrorNumber(c, err, 1301, http.StatusInternalServerError)
		return
	}

	jsonData := updatePlayerInviteRequest{}
	err = c.ShouldBindJSON(&jsonData)
	if err != nil {
		responseErrorNumber(c, err, 1301, http.StatusInternalServerError)
		return
	}

	if jsonData.InviteID == 0 {
		responseErrorNumber(c, nil, 1302, http.StatusBadRequest)
		return
	}

	invite, err := model.GetPlayerInviteToTeamById(jsonData.InviteID, model.TISended)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	if invite.ID == 0 {
		responseErrorNumber(c, nil, 1303, http.StatusBadRequest)
		return
	}

	if invite.Status == model.TICancel {
		responseErrorNumber(c, nil, 1304, http.StatusBadRequest)
		return
	}

	if invite.Status == model.TIRejected || invite.Status == model.TIAccepted {
		responseErrorNumber(c, nil, 1305, http.StatusBadRequest)
		return
	}

	if !strings.EqualFold(invite.Email, user.Email) {
		responseErrorNumber(c, nil, 1306, http.StatusBadRequest)
		return
	}

	if model.TeamInviteStatus(jsonData.Status) == model.TIAccepted {
		invite, err = model.AcceptedPlayerInviteToTeam(invite, userId)
		if err != nil {
			responseErrorNumber(c, err, 500, http.StatusInternalServerError)
			return
		}
	} else {
		invite.Status = model.TeamInviteStatus(jsonData.Status)
		invite, err = model.UpdatePlayerInviteToTeam(invite)
		if err != nil {
			responseErrorNumber(c, err, 500, http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, responseSuccess{
		Success: true,
	})
}
