package api

import (
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"

	"github.com/gin-gonic/gin"
)

type getOrganizationDataResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Address string `json:"addess"`
	UserID  uint   `json:"-"`
}
type getOrganizationResponse struct {
	Success bool                          `json:"success"`
	Data    []getOrganizationDataResponse `json:"data"`
}

// @Summary organization data
// @Schemes
// @Description get organization data
// @Tags profile organization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} getOrganizationResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/organization [get]
func GetOrganization(c *gin.Context) {
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

	organization, err := model.GetOrganizationByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1004,
			Message: GetMessageErr(1004),
		})
		return
	}
	result := []getOrganizationDataResponse{}

	if organization.ID > 0 {
		result = append(result, getOrganizationDataResponse{
			ID:      organization.ID,
			Title:   organization.Title,
			Address: organization.Address,
		})
	}

	c.JSON(http.StatusOK, getOrganizationResponse{
		Success: true,
		Data:    result,
	})
}

type createOrganizationRequest struct {
	Title   string `json:"title"`
	Address string `json:"addess"`
}

type createOrganizationResponse struct {
	Success bool                        `json:"success"`
	Data    getOrganizationDataResponse `json:"data"`
}

// @Summary create organization
// @Schemes
// @Description create organization
// @Tags profile organization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body createOrganizationRequest true "body"
// @Success 201 {object} createOrganizationResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/organization [post]
func CreateOrganization(c *gin.Context) {
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

	o, err := model.GetOrganizationByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1004,
			Message: GetMessageErr(1004),
		})
		return
	}
	if o.ID > 0 {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1007,
			Message: GetMessageErr(1007),
		})
		return
	}

	jsonData := createOrganizationRequest{}
	err = c.ShouldBindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1001,
			Message: GetMessageErr(1001),
		})
		return
	}

	// team, err := model.GetTeamByUser(userId)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, responseError{
	// 		Success: false,
	// 		Error:   1005,
	// 		Message: GetMessageErr(1005),
	// 	})
	// 	return
	// }
	// if team.ID > 0 {
	// 	c.JSON(http.StatusInternalServerError, responseError{
	// 		Success: false,
	// 		Error:   1006,
	// 		Message: GetMessageErr(1006),
	// 	})
	// 	return
	// }

	organization, err := model.CreateOrganization(jsonData.Title, userId, jsonData.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1002,
			Message: GetMessageErr(1002),
		})
		return
	}

	c.JSON(http.StatusCreated, createOrganizationResponse{
		Success: true,
		Data: getOrganizationDataResponse{
			ID:      organization.ID,
			Title:   organization.Title,
			Address: organization.Address,
		},
	})
}

type updateOrganizationRequest struct {
	Title   string `json:"title"`
	Address string `json:"addess"`
}

type updateOrganizationResponse struct {
	Success bool                        `json:"success"`
	Data    getOrganizationDataResponse `json:"data"`
}

// @Summary update organization
// @Schemes
// @Description update organization
// @Tags profile organization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param params body updateOrganizationRequest true "Organization"
// @Success 200 {object} updateOrganizationResponse
// @Failure 401 {object} responseError
// @Failure 401 {object} responseError
// @Failure 404 {object} responseError
// @Failure 500 {object} responseError
// @Router /profile/organization [put]
func UpdateOrganization(c *gin.Context) {
	session := sessions.New(c)

	jsonData := createOrganizationRequest{}
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1001,
			Message: GetMessageErr(1001),
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

	organization, err := model.GetOrganizationByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1004,
			Message: GetMessageErr(1004),
		})
		return
	}
	organization.Address = jsonData.Address
	organization.Title = jsonData.Title

	if organization.UserID != userId {
		c.JSON(http.StatusForbidden, responseError{
			Success: false,
			Error:   1008,
			Message: GetMessageErr(1008),
		})
		return
	}

	organization, err = model.UpdateOrganization(organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   1003,
			Message: GetMessageErr(1003),
		})
		return
	}

	c.JSON(http.StatusOK, updateOrganizationResponse{
		Success: true,
		Data: getOrganizationDataResponse{
			ID:      organization.ID,
			Title:   organization.Title,
			Address: organization.Address,
		},
	})
}
