package api

import (
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"

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
// @Success 200 {object} getTeamResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/team [get]
func GetTeam(c *gin.Context) {
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
			Error:   1205,
			Message: GetMessageErr(1205),
		})
		return
	}

	result := []getTeamDataResponse{}

	if t.ID > 0 {
		result = append(result, getTeamDataResponse{
			ID:      t.ID,
			Title:   t.Title,
			DGameID: t.DGameID,
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
