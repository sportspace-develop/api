package api

import (
	"fmt"
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"
	"sport-space-api/tools/email"
	"strconv"

	"github.com/gin-gonic/gin"
)

type getTeamDataResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	DGameID uint   `json:"game_type_id"`
}

type getTeamResponse struct {
	Success bool                  `json:"success"`
	Data    []getTeamDataResponse `json:"data"`
}

// type getTeamNullResponse struct {
// 	Success bool        `json:"success"`
// 	Data    interface{} `json:"data"`
// }

// @Summary team data
// @Schemes
// @Description get team data
// @Tags profile team
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param team_id query int false "Team id"
// @Success 200 {object} getTeamResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/team [get]
func GetTeam(c *gin.Context) {
	session := sessions.New(c)

	userId := session.GetUserId()
	if userId == 0 {
		responseErrorNumber(c, nil, 100, http.StatusUnauthorized)
		return
	}

	var err error
	var teamId int
	teamId, _ = strconv.Atoi(c.DefaultQuery("team_id", "0"))

	var teams []model.Team
	if teamId > 0 {
		t, err := model.GetTeamById(uint(teamId))
		if err != nil {
			responseErrorNumber(c, err, 500, http.StatusInternalServerError)
			return
		}
		if t.ID > 0 && t.UserID == userId {
			teams = append(teams, t)
		}
	} else {
		teams, err = model.GetTeamsByUser(userId)
		if err != nil {
			responseErrorNumber(c, err, 500, http.StatusInternalServerError)
			return
		}
	}

	result := []getTeamDataResponse{}

	for _, team := range teams {
		result = append(result, getTeamDataResponse{
			ID:      team.ID,
			Title:   team.Title,
			DGameID: team.DGameID,
		})
	}

	c.JSON(http.StatusOK, getTeamResponse{
		Success: true,
		Data:    result,
	})
}

type createTeamRequest struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	DGameID uint   `json:"game_type_id"`
}

type createTeamResponse struct {
	Success bool                `json:"success"`
	Data    getTeamDataResponse `json:"data"`
}

// @Summary create team
// @Schemes
// @Description create team
// @Tags profile team
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body createTeamRequest true "Team"
// @Success 200 {object} createTeamResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/team [post]
func CreateTeam(c *gin.Context) {
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

	t, err := model.GetTeamByUser(userId)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1200,
			Message: GetMessageErr(1200),
		})
		return
	}
	if t.ID > 0 {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1201,
			Message: GetMessageErr(1201),
		})
		return
	}

	jsonData := createTeamRequest{}
	err = c.ShouldBindJSON(&jsonData)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1203,
			Message: GetMessageErr(1203),
		})
		return
	}

	t, err = model.CreateTeam(jsonData.Title, jsonData.DGameID, userId)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1204,
			Message: GetMessageErr(1204),
		})
		return
	}

	c.JSON(http.StatusOK, createTeamResponse{
		Success: true,
		Data: getTeamDataResponse{
			ID:      t.ID,
			Title:   t.Title,
			DGameID: t.DGameID,
		},
	})
}

type updateTeamRequest struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	DGameID uint   `json:"game_type_id"`
}

type updateTeamResponse struct {
	Success bool                `json:"success"`
	Data    getTeamDataResponse `json:"data"`
}

// @Summary update team
// @Schemes
// @Description update team
// @Tags profile team
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body updateTeamRequest true "Team"
// @Success 200 {object} updateTeamResponse
// @Failure 401 {object} responseError
// @Failure 403 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/team [put]
func UpdateTeam(c *gin.Context) {
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

	jsonData := updateTeamRequest{}
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1203,
			Message: GetMessageErr(1203),
		})
		return
	}

	team, err := model.GetTeamById(jsonData.ID)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1205,
			Message: GetMessageErr(1205),
		})
		return
	}

	team.Title = jsonData.Title
	team.DGameID = jsonData.DGameID

	if team.UserID != userId {
		c.JSON(http.StatusForbidden, responseError{
			Success: false,
			Error:   1206,
			Message: GetMessageErr(1206),
		})
		return
	}

	result, err := model.UpdateTeam(team)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1206,
			Message: GetMessageErr(1206),
		})
		return
	}

	c.JSON(http.StatusOK, updateTeamResponse{
		Success: true,
		Data: getTeamDataResponse{
			ID:      result.ID,
			Title:   result.Title,
			DGameID: result.DGameID,
		},
	})
}

type getInviteToTeamData struct {
	Status string `json:"status"`
	Email  string `json:"email"`
}

type getInviteToTeamResponse struct {
	Success      bool                  `json:"success"`
	InviteStatus map[string]string     `json:"invite_status"`
	Data         []getInviteToTeamData `json:"data"`
}

// @Summary create invite to team
// @Schemes
// @Description create invite to team
// @Tags profile team
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} getInviteToTeamResponse
// @Failure 400 {object} responseError
// @Failure 401 {object} responseError
// @Failure 400 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/team/invite [get]
func GetInviteToTeam(c *gin.Context) {
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

	team, err := model.GetTeamByUser(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	if team.ID == 0 {
		responseErrorNumber(c, nil, 1300, http.StatusBadRequest)
		return
	}

	var result []getInviteToTeamData
	invites, err := model.GetInvitesToTeamByTeam(team.ID)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	for _, invite := range invites {
		result = append(result, getInviteToTeamData{
			Status: invite.Status.ToString(),
			Email:  invite.Email,
		})
	}

	c.JSON(http.StatusOK, getInviteToTeamResponse{
		Success:      true,
		InviteStatus: model.InviteStatus,
		Data:         result,
	})
}

type createInviteToTeamRequest struct {
	TeamID uint          `json:"team_id"`
	Email  []email.Email `json:"email"`
}

// @Summary create invite to team
// @Schemes
// @Description create invite to team
// @Tags profile team
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body createInviteToTeamRequest true "emails for invite"
// @Success 200 {object} responseSuccess
// @Failure 400 {object} responseError
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/team/invite [post]
func CreateInviteToTeam(c *gin.Context) {
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

	team, err := model.GetTeamByUser(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	if team.ID == 0 {
		c.JSON(http.StatusBadRequest, responseError{
			Success: false,
			Error:   1300,
			Message: GetMessageErr(1300),
		})
		return
	}

	jData := createInviteToTeamRequest{}
	err = c.ShouldBindJSON(&jData)
	if err != nil {
		responseErrorNumber(c, err, 1301, http.StatusInternalServerError)
		return
	}

	var invites []model.TeamInvite

	for _, inv := range jData.Email {
		invites = append(invites, model.TeamInvite{
			Email: string(inv),
			// Code:   tools.RandNumRunes(6),
			Status: model.TIWait,
			TeamID: team.ID,
		})
	}

	invites, err = model.CreateOrUpdateInvitesToTeam(invites, team.ID)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	for i, invite := range invites {
		email.AddMail(invite.Email, "invite to team", fmt.Sprintf("Invite to team <b>%s</b>", team.Title))
		invites[i].Status = model.TISended
	}

	_, err = model.UpdateInvitesToTeam(invites)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, responseSuccess{
		Success: true,
	})
}

type updateInviteToTeamRequest struct {
	TeamID uint                   `json:"team_id"`
	Email  []email.Email          `json:"email"`
	Status model.TeamInviteStatus `json:"status"`
}

// @Summary update invite to team
// @Schemes
// @Description update invite to team
// @Tags profile team
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body updateInviteToTeamRequest true "emails for invite"
// @Success 200 {object} responseSuccess
// @Failure 400 {object} responseError
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/team/invite [put]
func UpdateInviteToTeam(c *gin.Context) {
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

	team, err := model.GetTeamByUser(userId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	if team.ID == 0 {
		c.JSON(http.StatusBadRequest, responseError{
			Success: false,
			Error:   1300,
			Message: GetMessageErr(1300),
		})
		return
	}

	jData := updateInviteToTeamRequest{}
	err = c.ShouldBindJSON(&jData)
	if err != nil {
		responseErrorNumber(c, err, 1301, http.StatusInternalServerError)
		return
	}

	if jData.Status != model.TICancel && jData.Status != model.TIWait {
		responseErrorNumber(c, err, 1307, http.StatusBadRequest)
		return
	}

	var invites []model.TeamInvite

	for _, inv := range jData.Email {
		invites = append(invites, model.TeamInvite{
			Email: string(inv),
			// Code:   tools.RandNumRunes(6),
			Status: jData.Status,
			TeamID: team.ID,
		})
	}
	fmt.Println(invites)
	invites, err = model.CreateOrUpdateInvitesToTeam(invites, team.ID)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	if jData.Status == model.TIWait {
		for i, invite := range invites {
			email.AddMail(invite.Email, "invite to team", fmt.Sprintf("Invite to team <b>%s</b>", team.Title))
			invites[i].Status = model.TISended
		}

		_, err = model.UpdateInvitesToTeam(invites)
		if err != nil {
			responseErrorNumber(c, err, 500, http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, responseSuccess{
		Success: true,
	})
}
