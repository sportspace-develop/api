package api

import (
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type getTournamentDataResponse struct {
	ID                    uint   `json:"id"`
	Title                 string `json:"title" exampl:"Tournament #1"`
	StartDate             string `json:"start_date" example:"2023-12-31 15:04:05"`
	EndDate               string `json:"end_date" example:"2023-12-31 15:04:05"`
	StartRegistrationDate string `json:"start_registration_date" example:"2023-12-31 15:04:05"`
	EndRegistrationDate   string `json:"end_registration_date" example:"2023-12-31 15:04:05"`
	DGameID               uint   `json:"game_type_id" example:"1"` // вид спорта
	UserID                uint   `json:"-"`
	IsTeam                bool   `json:"is_team"` // команды или индивидульно
}

type getTournamentResponse struct {
	responseSuccess
	Data []getTournamentDataResponse `json:"data"`
}

// @Summary tournament data
// @Schemes
// @Description get tournament data
// @Tags profile organization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param tournament_id query int false "Tournaments"
// @Success 200 {object} getTournamentResponse
// @Failure 401 {object} responseError
// @Failure 404 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/tournament [get]
func GetTournament(c *gin.Context) {
	session := sessions.New(c)

	var tournamentId int
	tournamentId, _ = strconv.Atoi(c.DefaultQuery("tournament_id", "0"))

	response := getTournamentResponse{}
	response.Data = []getTournamentDataResponse{}
	if tournamentId > 0 {
		tData, err := model.GetTournamentById(uint(tournamentId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, responseError{
				Success: false,
				Error:   701,
				Message: GetMessageErr(701),
			})
			return
		}

		response.Success = true
		if tData.ID > 0 {
			response.Data = append(response.Data, getTournamentDataResponse{
				ID:                    tData.ID,
				Title:                 tData.Title,
				StartDate:             tData.StartDate.Format(time.DateTime),
				EndDate:               tData.EndDate.Format(time.DateTime),
				StartRegistrationDate: tData.StartRegistrationDate.Format(time.DateTime),
				EndRegistrationDate:   tData.EndRegistrationDate.Format(time.DateTime),
				IsTeam:                tData.IsTeam,
				DGameID:               tData.DGameID,
			})
		}
	} else {
		userId := session.GetUserId()
		if userId == 0 {
			c.JSON(http.StatusUnauthorized, responseError{
				Success: false,
				Error:   100,
				Message: GetMessageErr(100),
			})
			return
		}
		tData, err := model.GetTournamentsByUser(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responseError{
				Success: false,
				Error:   701,
				Message: GetMessageErr(701),
			})
			return
		}
		for _, row := range tData {
			response.Data = append(response.Data, getTournamentDataResponse{
				ID:                    row.ID,
				Title:                 row.Title,
				StartDate:             row.StartDate.Format(time.DateTime),
				EndDate:               row.EndDate.Format(time.DateTime),
				StartRegistrationDate: row.StartRegistrationDate.Format(time.DateTime),
				EndRegistrationDate:   row.EndRegistrationDate.Format(time.DateTime),
				IsTeam:                row.IsTeam,
				DGameID:               row.DGameID,
			})
		}
		response.Success = true
	}
	c.JSON(http.StatusOK, response)
}

type createTournamentRequest struct {
	Title                 string `json:"title" exampl:"Tournament #1"`
	StartDate             string `json:"start_date" example:"2023-12-31 00:00:00"`
	EndDate               string `json:"end_date" example:"2023-12-31 00:00:00"`
	StartRegistrationDate string `json:"start_registration_date" example:"2023-12-31 00:00:00"`
	EndRegistrationDate   string `json:"end_registration_date" example:"2023-12-31 00:00:00"`
	DGameID               uint   `json:"game_type_id" example:"1"` // вид спорта
	UserID                uint   `json:"-"`
	IsTeam                bool   `json:"is_team"` // команды или индивидульно
}

type createTournamentResponse struct {
	Success bool                      `json:"success"`
	Data    getTournamentDataResponse `json:"data"`
}

// @Summary create tournament
// @Schemes
// @Description create tournament
// @Tags profile organization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body createTournamentRequest true "Tournament"
// @Success 200 {object} createTournamentResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/tournament [post]
func CreateTournament(c *gin.Context) {
	session := sessions.New(c)

	var jsonData createTournamentRequest
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: MessageErr[500],
		})
		return
	}

	userId := session.GetUserId()
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, responseError{
			Success: false,
			Error:   100,
			Message: GetMessageErr(100),
		})
		return
	}

	startDate, err := time.Parse(time.DateTime, jsonData.StartDate)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   704,
			Message: GetMessageErr(702),
		})
		return
	}
	endDate, err := time.Parse(time.DateTime, jsonData.EndDate)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   705,
			Message: GetMessageErr(702),
		})
		return
	}
	startRegistrationDate, err := time.Parse(time.DateTime, jsonData.StartRegistrationDate)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   706,
			Message: GetMessageErr(702),
		})
		return
	}
	endRegistrationDate, err := time.Parse(time.DateTime, jsonData.EndRegistrationDate)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   707,
			Message: GetMessageErr(702),
		})
		return
	}

	t, err := model.CreateTournament(jsonData.Title, uint(userId), startDate, endDate, startRegistrationDate, endRegistrationDate, jsonData.IsTeam, jsonData.DGameID)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   703,
			Message: GetMessageErr(703),
		})
		return
	}

	c.JSON(http.StatusOK, createTournamentResponse{
		Success: true,
		Data: getTournamentDataResponse{
			ID:                    t.ID,
			Title:                 t.Title,
			StartDate:             t.StartDate.Format(time.DateTime),
			EndDate:               t.EndDate.Format(time.DateTime),
			StartRegistrationDate: t.StartRegistrationDate.Format(time.DateTime),
			EndRegistrationDate:   t.EndRegistrationDate.Format(time.DateTime),
			DGameID:               t.DGameID,
			UserID:                t.UserID,
			IsTeam:                t.IsTeam,
		},
	})
}

type updateTournamentRequest struct {
	ID                    uint   `json:"id"`
	Title                 string `json:"title" exampl:"Tournament #1"`
	StartDate             string `json:"start_date" example:"2023-12-31 00:00:00"`
	EndDate               string `json:"end_date" example:"2023-12-31 00:00:00"`
	StartRegistrationDate string `json:"start_registration_date" example:"2023-12-31 00:00:00"`
	EndRegistrationDate   string `json:"end_registration_date" example:"2023-12-31 00:00:00"`
	DGameID               uint   `json:"game_type_id" example:"1"` // вид спорта
	UserID                uint   `json:"-"`
	IsTeam                bool   `json:"is_team"` // команды или индивидульно
}
type updateTournamentResponse struct {
	Success bool                      `json:"success"`
	Data    getTournamentDataResponse `json:"data"`
}

// @Summary edit tournament
// @Schemes
// @Description edit tournament
// @Tags profile organization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body updateTournamentRequest true "Tournaments"
// @Success 200 {object} updateTournamentResponse
// @Failure 401 {object} responseError
// @Failure 403 {object} responseError
// @Failure 404 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/tournament [put]
func UpdateTournament(c *gin.Context) {

	session := sessions.New(c)

	var jsonData updateTournamentRequest
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   700,
			Message: GetMessageErr(700),
		})
		return
	}

	t, err := model.GetTournamentById(jsonData.ID)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   700,
			Message: GetMessageErr(700),
		})
		return
	}

	userId := session.GetUserId()
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, responseError{
			Success: false,
			Error:   100,
			Message: GetMessageErr(100),
		})
		return
	}

	if t.UserID != userId {
		c.JSON(http.StatusForbidden, responseError{
			Success: false,
			Error:   702,
			Message: GetMessageErr(702),
		})
		return
	}

	startDate, err := time.Parse(time.DateTime, jsonData.StartDate)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   704,
			Message: GetMessageErr(704),
		})
		return
	}
	t.StartDate = startDate

	endDate, err := time.Parse(time.DateTime, jsonData.EndDate)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   705,
			Message: GetMessageErr(705),
		})
		return
	}
	t.EndDate = endDate

	startRegistrationDate, err := time.Parse(time.DateTime, jsonData.StartRegistrationDate)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   706,
			Message: GetMessageErr(706),
		})
		return
	}
	t.StartRegistrationDate = startRegistrationDate

	endRegistrationDate, err := time.Parse(time.DateTime, jsonData.EndRegistrationDate)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   707,
			Message: GetMessageErr(707),
		})
		return
	}
	t.EndRegistrationDate = endRegistrationDate

	t, err = model.UpdateTournament(t)
	if err != nil {
		log.ERROR(err.Error())
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   703,
			Message: GetMessageErr(703),
		})
		return
	}

	c.JSON(http.StatusOK, updateTournamentResponse{
		Success: true,
		Data: getTournamentDataResponse{
			ID:                    t.ID,
			Title:                 t.Title,
			StartDate:             t.StartDate.Format(time.DateTime),
			EndDate:               t.EndDate.Format(time.DateTime),
			StartRegistrationDate: t.StartRegistrationDate.Format(time.DateTime),
			EndRegistrationDate:   t.EndRegistrationDate.Format(time.DateTime),
			DGameID:               t.DGameID,
			UserID:                t.UserID,
			IsTeam:                t.IsTeam,
		},
	})
}
