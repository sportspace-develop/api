package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getUserResponse struct {
}

// @Summary user data
// @Schemes
// @Description get user data
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} getUserResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /user [get]
func GetUser(c *gin.Context) {

	c.JSON(http.StatusOK, responseSuccess{
		Success: true,
	})
}
