package api

const (
	DATE_FORMAT = "02.01.2006 15:04:05"
)

type responseError struct {
	Success      bool   `json:"success" swaggertype:"boolean" example:"false"`
	Error        int16  `json:"error"`
	Message      string `json:"message"`
	ErrorMessage string `json:"errorMessage"`
}

type responseSuccess struct {
	Success bool `json:"success" swaggertype:"boolean" example:"true"`
}
