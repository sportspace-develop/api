package api

import (
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type getPlayerDataResponse struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
	BDay       string `json:"b_day" example:"2000-12-31`
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

	player, err := model.GetPlayerByUserId(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, getPlayerResponse{
		Success: true,
		Data: getPlayerDataResponse{
			FirstName:  player.FirstName,
			SecondName: player.SecondName,
			LastName:   player.LastName,
			BDay:       player.BDay.Time.Format(time.DateOnly),
		},
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
