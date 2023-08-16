package api

import "sport-space-api/logger"

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
	Error   int16  `json:"error"`
	Message string `json:"message"`
}

type responseSuccess struct {
	Success bool `json:"success" swaggertype:"boolean" example:"true"`
}
