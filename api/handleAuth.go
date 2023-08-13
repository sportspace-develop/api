package api

import (
	"net/http"
	"sport-space-api/model"
	"sport-space-api/tools"
	"sport-space-api/tools/email"
	"sport-space-api/tools/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email string `json:"email" swaggertype:"string" example:"test@test.ru"`
}

type loginResponse struct {
	Success bool `json:"success"`
}

// @Summary authorization
// @Schemes
// @Description send code to email
// @Tags auth
// @Accept json
// @Produce json
// @Param email body loginRequest true "User email"
// @Success 200 {object} loginResponse
// @Failure 400 {object} responseError
// @Failure 500 {object} responseError
// @Router /auth/otp [post]
func GetAuthCode(c *gin.Context) {
	var jsonData loginRequest
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: MessageErr[500],
		})
		log.ERROR(err.Error())
		return
	}

	if jsonData.Email == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Success: false,
			Error:   12,
			Message: MessageErr[12],
		})
		return
	}

	user, err := model.FindUserByEmail(jsonData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}
	if user.ID == 0 {
		user, err = model.Registration(jsonData.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responseError{
				Success: false,
				Error:   500,
				Message: GetMessageErr(500),
			})
			log.ERROR(err.Error())
			return
		}
	}

	code := tools.RandNumRunes(6)

	mCode, err := model.FindCodeNotActivatedByUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}

	if mCode.ID != 0 {
		mCode.Code = code
		mCode.AttemptCount += 1
		mCode.AttemptDate = time.Now().UTC()
		mCode.ExpiresIn = time.Now().UTC().Add(time.Duration(time.Minute * 60))
		model.SaveCodeNotActivatedByUser(mCode)

		if sendOk, sendErr := email.SendCodeToEmail(jsonData.Email, code); sendErr != nil || !sendOk {
			c.JSON(http.StatusInternalServerError, responseError{
				Success: false,
				Error:   500,
				Message: GetMessageErr(500),
			})
			log.ERROR(err.Error())
			return
		}
	} else if ok, err := model.CreateUserAuthCode(user, code); err == nil && ok {
		if sendOk, sendErr := email.SendCodeToEmail(jsonData.Email, code); sendErr != nil || !sendOk {
			c.JSON(http.StatusInternalServerError, responseError{
				Success: false,
				Error:   500,
				Message: GetMessageErr(500),
			})
			log.ERROR(err.Error())
			return
		}

	} else {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}

	c.JSON(http.StatusOK, loginResponse{Success: true})
}

type authorizeRequest struct {
	Email email.Email `json:"email" swaggertype:"string" example:"test@test.ru"`
	Code  string      `json:"code" swaggertype:"string" example:"123456"`
}

type authorizeResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token" swaggertype:"string" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.-5myAJwbMszwt7_iPciBQgICdujy20zKOZOUTXu9KyY"`
	RefreshToken string `json:"refresh_token" swaggertype:"string" example:"cyYTkJzAjEAgcaIIUPeZvyLpZHVuBIArVXqpInHLrbvXzgofSWKWlbZflPUToIctnWJoJInIqfDVLTIOeBGtJMRnlhseRgpHlPxh"`
	ExpiresIn    string `json:"expires_in" swaggertype:"string" example:"2006-01-02 15:04:05"`
}

// @Summary authorization
// @Schemes
// @Description user authorization
// @Tags auth
// @Accept json
// @Produce json
// @Param params body authorizeRequest true "User email"
// @Success 200 {object} authorizeResponse
// @Failure 400 {object} responseError
// @Failure 500 {object} responseError
// @Router /auth/login [post]
func Authorize(c *gin.Context) {
	var jsonData authorizeRequest
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: MessageErr[500],
		})
		return
	}

	if !jsonData.Email.IsValid() {
		c.JSON(http.StatusBadRequest, responseError{
			Success: false,
			Error:   12,
			Message: GetMessageErr(12),
		})
		return
	}

	user, err := model.FindUserByEmail(string(jsonData.Email))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}

	authCode, err := model.FindCodeNotActivatedByUserCode(user, jsonData.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}

	if authCode.ID == 0 {
		c.JSON(http.StatusBadRequest, responseError{
			Success: false,
			Error:   11,
			Message: GetMessageErr(11),
		})
		return
	}

	_, err = model.ActivateUserAuthCode(authCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Success: false,
			Error:   17,
			Message: GetMessageErr(17),
		})
		log.ERROR(err.Error())
		return
	}

	tokenString, err := jwt.New(map[jwt.Fields]interface{}{
		jwt.USER_ID: user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   13,
			Message: GetMessageErr(13),
		})
		log.ERROR(err.Error())
		return
	}

	refreshToken := tools.RandStringRunes(100)

	newSess, err := model.NewSession(user, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}

	c.JSON(http.StatusOK, authorizeResponse{
		Success:      true,
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    newSess.ExpiresIn.Format(time.DateTime),
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" swaggertype:"string" example:"qweqwe231e2qeqae"`
}

type refreshResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token" swaggertype:"string" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.-5myAJwbMszwt7_iPciBQgICdujy20zKOZOUTXu9KyY"`
	RefreshToken string `json:"refresh_token" swaggertype:"string" example:"cyYTkJzAjEAgcaIIUPeZvyLpZHVuBIArVXqpInHLrbvXzgofSWKWlbZflPUToIctnWJoJInIqfDVLTIOeBGtJMRnlhseRgpHlPxh"`
	ExpiresIn    string `json:"expires_in" swaggertype:"string" example:"2006-01-02 15:04:05"`
}

// @Summary authorization
// @Schemes
// @Description refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param params body refreshRequest true "Refresh token"
// @Success 200 {object} refreshResponse
// @Failure 400 {object} responseError
// @Failure 500 {object} responseError
// @Router /auth/refresh [post]
func Refresh(c *gin.Context) {
	jData := refreshRequest{}
	err := c.ShouldBindJSON(&jData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: MessageErr[500],
		})
		log.ERROR(err.Error())
		return
	}

	//проверяем refresh-token
	if jData.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Success: false,
			Error:   20,
			Message: GetMessageErr(20),
		})
		return
	}

	sess, err := model.FindSession(jData.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   13,
			Message: GetMessageErr(13),
		})
		log.ERROR(err.Error())
		return
	}
	if sess.ID != 0 {
		_, err := model.DeleteSession(sess.RefreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, responseError{
				Success: false,
				Error:   500,
				Message: GetMessageErr(500),
			})
			log.ERROR(err.Error())
			return
		}
	}
	if sess.ID == 0 {
		c.JSON(http.StatusBadRequest, responseError{
			Success: false,
			Error:   21,
			Message: GetMessageErr(21),
		})
		return
	}

	user, err := model.FindUserById(sess.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, responseError{
			Success: false,
			Error:   21,
			Message: GetMessageErr(21),
		})
		return
	}

	tokenString, err := jwt.New(map[jwt.Fields]interface{}{
		jwt.USER_ID: user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   13,
			Message: GetMessageErr(13),
		})
		log.ERROR(err.Error())
		return
	}
	refreshToken := tools.RandStringRunes(100)
	newSess, err := model.NewSession(user, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}

	c.JSON(http.StatusOK, refreshResponse{
		Success:      true,
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    newSess.ExpiresIn.Format(time.DateTime),
	})
}
