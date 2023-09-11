package api

import (
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"

	"github.com/gin-gonic/gin"
)

type getProfileResponse struct {
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
