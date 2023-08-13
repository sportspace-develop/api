package api

import (
	"encoding/json"
	"net/http"
	"sport-space-api/tools/jwt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRequiredMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		var h_auth string
		if len(c.Request.Header["Authorization"]) > 0 {
			h_auth = c.Request.Header["Authorization"][0]
		}

		h_auth_splited := strings.Split(h_auth, " ")
		if len(h_auth_splited) > 1 {
			h_auth = h_auth_splited[1]
		}

		token, err := jwt.Check(h_auth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, responseError{
				Success: false,
				Error:   14,
				Message: GetMessageErr(14),
			})
			log.ERROR(err.Error())
			c.Abort()
			return
		}

		isExpired, err := jwt.IsExpired(token)
		if !token.Valid || isExpired || err != nil {
			c.JSON(http.StatusUnauthorized, responseError{
				Success: false,
				Error:   401,
				Message: MessageErr[401],
			})
			log.ERROR(err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		duration := time.Since(start).Seconds()
		fields := map[string]interface{}{
			"client_ip":            c.ClientIP(),
			"duration":             duration,
			"method":               c.Request.Method,
			"header_authorization": c.Request.Header["Authorization"],
			"path":                 c.Request.URL.Path,
			"request":              c.Request.URL.RawQuery,
			"status":               c.Writer.Status(),
			// "user_id":    0, //util.GetUserID(c),
			"referrer":   c.Request.Referer(),
			"request_id": c.Writer.Header().Get("Request-Id"),
		}
		if c.Request.Method != "GET" {
			var jsonDataBytes interface{} = make(map[string]interface{})
			c.Copy().ShouldBindJSON(&jsonDataBytes)
			fields["request_body"] = jsonDataBytes
		}

		marshaled, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			log.ERROR("marshaling error: %s", err.Error())
		}

		log.INFO(string(marshaled))
	}
}
