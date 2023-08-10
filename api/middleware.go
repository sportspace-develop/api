package api

import (
	"net/http"
	"sport-space-api/config"
	"sport-space-api/tools/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {

	return func(c *gin.Context) {
		h_auth := c.Request.Header["Authorization"][0]

		h_auth_splited := strings.Split(h_auth, " ")
		if len(h_auth_splited) > 1 {
			h_auth = h_auth_splited[1]
		}

		token, err := jwt.Check(h_auth, []byte(config.App.JWTSecret))
		if err != nil {
			c.JSON(http.StatusUnauthorized, responseError{
				Success:      false,
				Error:        14,
				Message:      GetMessageErr(14),
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		isExpired, err := jwt.IsExpired(token)
		if !token.Valid || isExpired || err != nil {
			c.JSON(http.StatusUnauthorized, responseError{
				Success:      false,
				Error:        401,
				Message:      MessageErr[401],
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
