package api

import (
	"encoding/json"
	"net/http"
	sessions "sport-space-api/session"
	"sport-space-api/tools/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRequiredMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		session := sessions.New(c)

		token, err := jwt.DefaultGin(c)
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

		isExpired, err := token.IsExpired()
		if !token.Valid() || isExpired || err != nil {
			c.JSON(http.StatusUnauthorized, responseError{
				Success: false,
				Error:   401,
				Message: GetMessageErr(401),
			})
			if err != nil {
				log.ERROR(err.Error())
			}
			c.Abort()
			return
		}
		tUserId, _ := token.GetUserId()
		if tUserId != session.GetUserId() && tUserId > 0 {
			session.Clear()
			session.SetUserId(tUserId)
			session.Save()
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

		var request_body string
		// if c.Request.Method != "GET" {

		// 	jsonData, _ := ioutil.ReadAll(c.Request.Body)
		// 	request_body = string(jsonData)
		// }

		// Process Request
		c.Next()

		token, err := jwt.DefaultGin(c)
		var userId uint
		if err == nil && token.Valid() {
			userId, _ = token.GetUserId()
		}

		duration := time.Since(start).Seconds()
		fields := map[string]interface{}{
			"client_ip":            c.ClientIP(),
			"duration":             duration,
			"method":               c.Request.Method,
			"header_authorization": c.Request.Header["Authorization"],
			"path":                 c.Request.URL.Path,
			"request":              c.Request.URL.RawQuery,
			"status":               c.Writer.Status(),
			"user_id":              userId,
			"referrer":             c.Request.Referer(),
			"request_id":           c.Writer.Header().Get("Request-Id"),
			"request_body":         request_body,
		}

		marshaled, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			log.ERROR("marshaling error: %s", err.Error())
		}

		log.INFO(string(marshaled))
	}
}
