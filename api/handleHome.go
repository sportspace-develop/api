package api

import (
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type getAllOrganizationDataResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Address string `json:"addess"`
	UserID  uint   `json:"-"`
}
type getAllOrganizationResponse struct {
	Success bool                             `json:"success"`
	Data    []getAllOrganizationDataResponse `json:"data"`
}

// @Summary info for all users
// @Schemes
// @Description get organization data
// @Tags Home
// @Accept json
// @Produce json
// @Success 200 {object} getAllOrganizationResponse
// @Failure 500 {object} responseError
// @Router /organization [get]
func GetAllOrganization(c *gin.Context) {
	// session := sessions.New(c)

	// userId := session.GetUserId()
	// if userId == 0 {
	// 	c.JSON(http.StatusUnauthorized, responseError{
	// 		Success: false,
	// 		Error:   100,
	// 		Message: GetMessageErr(100),
	// 	})
	// 	return
	// }

	organizations, err := model.GetOrganizations()
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	result := []getAllOrganizationDataResponse{}

	for _, organization := range organizations {
		result = append(result, getAllOrganizationDataResponse{
			ID:      organization.ID,
			Title:   organization.Title,
			Address: organization.Address,
		})
	}

	c.JSON(http.StatusOK, getAllOrganizationResponse{
		Success: true,
		Data:    result,
	})
}

type getTournamentsDataResponse struct {
	ID                    uint   `json:"id"`
	Title                 string `json:"title" exampl:"Tournament #1"`
	StartDate             string `json:"start_date" example:"2023-12-31 15:04:05"`
	EndDate               string `json:"end_date" example:"2023-12-31 15:04:05"`
	StartRegistrationDate string `json:"start_registration_date" example:"2023-12-31 15:04:05"`
	EndRegistrationDate   string `json:"end_registration_date" example:"2023-12-31 15:04:05"`
	DGameID               uint   `json:"game_type_id" example:"1"` // вид спорта
	UserID                uint   `json:"-"`
	IsTeam                bool   `json:"is_team"` // команды или индивидульно
	OrganizationID        uint   `json:"organization_id"`
}

type getTournamentsResponse struct {
	responseSuccess
	Data []getTournamentsDataResponse `json:"data"`
}

// @Summary tournament data
// @Schemes
// @Description get tournament data
// @Tags Home
// @Accept json
// @Produce json
// @Param tournament_id query int false "Tournaments"
// @Success 200 {object} getTournamentsResponse
// @Failure 404 {object} responseError
// @Failure 500 {object} responseError
// @Router /tournament [get]
func GetTournaments(c *gin.Context) {

	var tournamentId int
	tournamentId, _ = strconv.Atoi(c.DefaultQuery("tournament_id", "0"))

	response := getTournamentsResponse{}
	response.Data = []getTournamentsDataResponse{}
	tournaments, err := model.GetTournaments(tournamentId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	response.Success = true
	for _, t := range tournaments {
		response.Data = append(response.Data, getTournamentsDataResponse{
			ID:                    t.ID,
			Title:                 t.Title,
			StartDate:             t.StartDate.Format(time.DateTime),
			EndDate:               t.EndDate.Format(time.DateTime),
			StartRegistrationDate: t.StartRegistrationDate.Format(time.DateTime),
			EndRegistrationDate:   t.EndRegistrationDate.Format(time.DateTime),
			IsTeam:                t.IsTeam,
			DGameID:               t.DGameID,
			OrganizationID:        t.OrganizationID,
		})
	}
	c.JSON(http.StatusOK, response)
}

type getTeamsDataResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	DGameID uint   `json:"game_type_id"`
}

type getTeamsResponse struct {
	Success bool                   `json:"success"`
	Data    []getTeamsDataResponse `json:"data"`
}

// type getTeamNullResponse struct {
// 	Success bool        `json:"success"`
// 	Data    interface{} `json:"data"`
// }

// @Summary team data
// @Schemes
// @Description get team data
// @Tags Home
// @Accept json
// @Produce json
// @Param team_id query int false "Team"
// @Success 200 {object} getTeamResponse
// @Failure 500 {object} responseError
// @Router /team [get]
func GetTeams(c *gin.Context) {

	var teamId int
	teamId, _ = strconv.Atoi(c.DefaultQuery("team_id", "0"))

	teams, err := model.GetTeams(teamId)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	result := []getTeamsDataResponse{}

	for _, team := range teams {
		result = append(result, getTeamsDataResponse{
			ID:      team.ID,
			Title:   team.Title,
			DGameID: team.DGameID,
		})
	}

	c.JSON(http.StatusOK, getTeamsResponse{
		Success: true,
		Data:    result,
	})
}

type getHomeResponse struct {
	Success       bool                             `json:"success"`
	UserID        uint                             `json:"user_Id"`
	GameTypes     []getProfileGameTypesResponse    `json:"game_types"`
	Organizations []getAllOrganizationDataResponse `json:"organizations"`
	Tournaments   []getTournamentsDataResponse     `json:"tournaments"`
	Teams         []getTeamsDataResponse           `json:"teams"`
}

// @Summary home data
// @Schemes
// @Description get home data
// @Tags Home
// @Accept json
// @Produce json
// @Success 200 {object} getHomeResponse
// @Failure 500 {object} responseError
// @Router /home [get]
func GetHome(c *gin.Context) {
	session := sessions.New(c)
	userId := session.GetUserId()

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

	organizations, err := model.GetOrganizations()
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	resultOrganizations := []getAllOrganizationDataResponse{}

	for _, organization := range organizations {
		resultOrganizations = append(resultOrganizations, getAllOrganizationDataResponse{
			ID:      organization.ID,
			Title:   organization.Title,
			Address: organization.Address,
		})
	}

	tournaments, err := model.GetTournaments[int](0)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}
	resultTournaments := []getTournamentsDataResponse{}
	for _, t := range tournaments {
		resultTournaments = append(resultTournaments, getTournamentsDataResponse{
			ID:                    t.ID,
			Title:                 t.Title,
			StartDate:             t.StartDate.Format(time.DateTime),
			EndDate:               t.EndDate.Format(time.DateTime),
			StartRegistrationDate: t.StartRegistrationDate.Format(time.DateTime),
			EndRegistrationDate:   t.EndRegistrationDate.Format(time.DateTime),
			IsTeam:                t.IsTeam,
			DGameID:               t.DGameID,
			OrganizationID:        t.OrganizationID,
		})
	}

	teams, err := model.GetTeams(0)
	if err != nil {
		responseErrorNumber(c, err, 500, http.StatusInternalServerError)
		return
	}

	resultTeams := []getTeamsDataResponse{}

	for _, team := range teams {
		resultTeams = append(resultTeams, getTeamsDataResponse{
			ID:      team.ID,
			Title:   team.Title,
			DGameID: team.DGameID,
		})
	}

	c.JSON(http.StatusOK, getHomeResponse{
		Success:       true,
		UserID:        userId,
		GameTypes:     gameTypesResult,
		Organizations: resultOrganizations,
		Tournaments:   resultTournaments,
		Teams:         resultTeams,
	})
}
