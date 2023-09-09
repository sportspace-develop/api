package jwt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	gJWT "github.com/golang-jwt/jwt/v5"
)

type Fields string

const (
	USER_ID Fields = "user_id"
)

var (
	Secret              []byte
	AccessTokenLongTime time.Duration
)

type tokenT struct {
	t *gJWT.Token
}

func (t *tokenT) String() (string, error) {
	return t.t.SignedString(Secret)
}

func Parse(token string) (*tokenT, error) {

	tStr, err := Check(token)

	return &tokenT{
		t: tStr,
	}, err
}

func New(jData map[Fields]interface{}) *tokenT {

	fmt.Println("access long", AccessTokenLongTime)
	expires := time.Now().UTC().Add(AccessTokenLongTime)
	jwtData := gJWT.MapClaims{
		"expires_in": expires.Format(time.DateTime),
	}

	for k, v := range jData {
		jwtData[string(k)] = v
	}

	token := gJWT.NewWithClaims(gJWT.SigningMethodHS256, jwtData)
	tkn := tokenT{
		t: token,
	}
	return &tkn
}

func (t *tokenT) IsExpired() (bool, error) {
	if claims, ok := t.t.Claims.(gJWT.MapClaims); ok && t.t.Valid {
		if _, ok := claims["expires_in"]; !ok {
			return true, fmt.Errorf("Token not valid")
		}

		expired_dt, err := time.Parse(time.DateTime, claims["expires_in"].(string))
		if err != nil {
			return true, fmt.Errorf("Not found attribute expires_in")
		}
		if expired_dt.Sub(time.Now().UTC()) < time.Duration(0) {
			return true, nil
		}
	}

	return false, nil

}

func (t *tokenT) Valid() bool {
	return t.t.Valid
}

func (t *tokenT) GetExpiresDateString() (string, error) {
	claims, ok := t.t.Claims.(gJWT.MapClaims)
	if !ok {
		return "", fmt.Errorf("Token not valid")
	}
	if _, ok := claims["expires_in"]; !ok {
		return "", fmt.Errorf("Token not valid")
	}

	return claims["expires_in"].(string), nil
}

func (t *tokenT) GetUserId() (uint, error) {
	claims, ok := t.t.Claims.(gJWT.MapClaims)
	if !ok {
		return 0, fmt.Errorf("Token not valid")
	}
	if _, ok := claims["expires_in"]; !ok {
		return 0, fmt.Errorf("Token not valid")
	}
	sUserId := claims["user_id"].(string)
	userId, err := strconv.Atoi(sUserId)
	if err != nil {
		return 0, nil
	}
	return uint(userId), nil
}

func Check(hToken string) (*gJWT.Token, error) {

	token, err := gJWT.Parse(hToken, func(token *gJWT.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*gJWT.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return Secret, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func DefaultGin(c *gin.Context) (*tokenT, error) {
	var h_auth string
	if len(c.Request.Header["Authorization"]) > 0 {
		h_auth = c.Request.Header["Authorization"][0]
	}

	h_auth = strings.ReplaceAll(h_auth, "Bearer ", "")
	token, err := Parse(h_auth)
	return token, err
}
