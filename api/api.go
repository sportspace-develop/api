package api

import (
	"net/http"
	"sport-space-api/logger"

	"github.com/gin-gonic/gin"
)

const (
	DATE_FORMAT = "02.01.2006 15:04:05"
)

var (
	log *logger.Logger
)

func Init() {
	log = logger.New("api")
	log.INFO("init api")
}

type responseError struct {
	Success bool   `json:"success" swaggertype:"boolean" example:"false"`
	Error   int    `json:"error"`
	Message string `json:"message"`
}

type responseSuccess struct {
	Success bool `json:"success" swaggertype:"boolean" example:"true"`
}

func responseErrorNumber(c *gin.Context, err error, errNum int, statusCode int) {
	if err != nil {
		log.ERROR(err.Error())
	}
	c.JSON(http.StatusInternalServerError, responseError{
		Success: false,
		Error:   errNum,
		Message: GetMessageErr(errNum),
	})
}
