package api

import (
	"fmt"
	"net/http"
	"sport-space-api/model"
	sessions "sport-space-api/session"
	"sport-space-api/tools"
	"sport-space-api/tools/email"
	"sport-space-api/tools/jwt"
	"sport-space-api/tools/password"
	"time"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email string `json:"email" swaggertype:"string" example:"test@test.ru"`
}

type loginResponse struct {
	Success bool `json:"success"`
}

// @Summary send to email one time password
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
		model.UpdateCodeNotActivatedByUser(mCode)
		email.AddMail(jsonData.Email, "Auth code", code)
	} else if ok, err := model.CreateUserAuthCode(user, code); err == nil && ok {
		email.AddMail(jsonData.Email, "Auth code", code)
	}

	c.JSON(http.StatusOK, loginResponse{Success: true})
}

type authAction string

const (
	ActionCode     = "code"
	ActionPassword = "password"
)

type authorizeRequest struct {
	Email    email.Email `json:"email" swaggertype:"string" example:"test@test.ru"`
	Action   authAction  `json:"action" swaggertype:"string" example:"code"`
	Code     string      `json:"code" swaggertype:"string" example:"123456"`
	Password string      `json:"password" swaggertype:"string" example:"password_string"`
}

type authorizeResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token" swaggertype:"string" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.-5myAJwbMszwt7_iPciBQgICdujy20zKOZOUTXu9KyY"`
	RefreshToken string `json:"refresh_token" swaggertype:"string" example:"cyYTkJzAjEAgcaIIUPeZvyLpZHVuBIArVXqpInHLrbvXzgofSWKWlbZflPUToIctnWJoJInIqfDVLTIOeBGtJMRnlhseRgpHlPxh"`
	ExpiresIn    string `json:"expires_in" swaggertype:"string" example:"2006-01-02 15:04:05"`
}

// @Summary authorization
// @Schemes
// @Description Авторизация по паролю или по коду из email, action = "code" | "password". Для code обязательное поле "code", для пароля обязательное поле "password"
// @Tags auth
// @Accept json
// @Produce json
// @Param params body authorizeRequest true "User email"
// @Success 200 {object} authorizeResponse
// @Failure 400 {object} responseError
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Router /auth/login [post]
func Authorize(c *gin.Context) {
	session := sessions.New(c)

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

	// if user.ID == 0 {
	// 	responseErrorNumber(c, nil, 19, http.StatusUnauthorized)
	// 	return
	// }

	if jsonData.Action == ActionCode {
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
			responseErrorNumber(c, nil, 11, http.StatusBadRequest)
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
	} else if jsonData.Action == ActionPassword {
		if !password.CheckPasswordHash(jsonData.Password, user.Password.String) {
			responseErrorNumber(c, nil, 19, http.StatusUnauthorized)
			return
		}

	} else {
		responseErrorNumber(c, nil, 9, http.StatusBadRequest)
		return
	}

	tkn := jwt.New(map[jwt.Fields]interface{}{
		jwt.USER_ID: fmt.Sprint(user.ID),
	})
	tokenString, err := tkn.String()
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

	_, err = model.NewSession(user, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}
	expiresIn, err := tkn.GetExpiresDateString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}
	userId, err := tkn.GetUserId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   18,
			Message: GetMessageErr(18),
		})
		log.ERROR(err.Error())
		return
	}
	session.Clear()
	session.SetUserId(userId)
	session.Set("refresh_token", refreshToken)
	session.Save()

	c.JSON(http.StatusOK, authorizeResponse{
		Success:      true,
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" swaggertype:"string" example:"cyYTkJzAjEAgcaIIUPeZvyLpZHVuBIArVXqpInHLrbvXzgofSWKWlbZflPUToIctnWJoJInIqfDVLTIOeBGtJMRnlhseRgpHlPxh"`
}

type refreshResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token" swaggertype:"string" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.-5myAJwbMszwt7_iPciBQgICdujy20zKOZOUTXu9KyY"`
	RefreshToken string `json:"refresh_token" swaggertype:"string" example:"cyYTkJzAjEAgcaIIUPeZvyLpZHVuBIArVXqpInHLrbvXzgofSWKWlbZflPUToIctnWJoJInIqfDVLTIOeBGtJMRnlhseRgpHlPxh"`
	ExpiresIn    string `json:"expires_in" swaggertype:"string" example:"2006-01-02 15:04:05"`
}

// @Summary refresh token
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
	session := sessions.New(c)

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

	tkn := jwt.New(map[jwt.Fields]interface{}{
		jwt.USER_ID: fmt.Sprint(user.ID),
	})
	tokenString, err := tkn.String()
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
	_, err = model.NewSession(user, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}

	expiresIn, err := tkn.GetExpiresDateString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}

	session.Set("refresh_token", refreshToken)
	session.Save()

	c.JSON(http.StatusOK, refreshResponse{
		Success:      true,
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	})
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token" swaggertype:"string" example:"cyYTkJzAjEAgcaIIUPeZvyLpZHVuBIArVXqpInHLrbvXzgofSWKWlbZflPUToIctnWJoJInIqfDVLTIOeBGtJMRnlhseRgpHlPxh"`
}

// @Summary logout
// @Schemes
// @Description logout
// @Tags auth
// @Accept json
// @Produce json
// @Param params body logoutRequest true "Refresh token"
// @Param Authorization header string true "Bearer JWT"
// @Success 401 {object} responseSuccess
// @Failure 500 {object} responseError
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	session := sessions.New(c)

	// refreshToken := session.Get("refresh_token").(string)

	var jsonData logoutRequest
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: MessageErr[500],
		})
		return
	}

	_, err = model.DeleteSession(jsonData.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Success: false,
			Error:   500,
			Message: GetMessageErr(500),
		})
		log.ERROR(err.Error())
		return
	}

	session.Clear()
	session.Save()
	c.JSON(http.StatusUnauthorized, responseSuccess{
		Success: true,
	})
}
