package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"sport-space/internal/adapter/models"
	"sport-space/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	errUnauthorize = errors.New("unauthorize")
)

var (
	cookieName = "token"
	cookieKey  = "UserID"
)

func (s *Server) middlewareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		s.log.Info(
			"Request",
			zap.String("uri", c.Request.RequestURI),
			zap.Duration("duration", time.Since(start)),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
		)
	}
}

func (s *Server) middlewareAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := s.checkAuth(c)
		if err != nil {
			if !errors.Is(err, errUnauthorize) {
				c.Writer.WriteHeader(http.StatusInternalServerError)
			} else {
				c.Writer.WriteHeader(http.StatusUnauthorized)
			}
			c.Abort()
		}

		if err == nil && userID == 0 {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Abort()
		}

		c.Next()
	}
}

func (s *Server) checkAuth(c *gin.Context) (userID uint, err error) {
	var ok bool
	var userIDS string
	cookieUserID, err := c.Request.Cookie(cookieName)
	if err != nil {
		return 0, fmt.Errorf("failed reade user cookie: %w %w", err, errUnauthorize)
	}

	jwtRest := jwt.New([]byte(s.secret))
	userIDS, ok, err = jwtRest.Verify(cookieUserID.Value, cookieKey)
	if err != nil {
		return 0, fmt.Errorf("failed verify token: %w %w", err, errUnauthorize)
	}

	if !ok {
		return 0, fmt.Errorf("unverify usercookie: %w", errUnauthorize)
	}

	userID64, err := strconv.ParseUint(userIDS, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("can't convert string userID to uint: %w", err)
	}

	return uint(userID64), nil
}

func (s *Server) checkUser(c *gin.Context) (user *models.User, statusCode int, err error) {
	userID, err := s.checkAuth(c)
	if err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("failed check cookie: %w", err)
	}
	user, err = s.sport.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed get user: %w", err)
	}
	if user.ID == 0 {
		return nil, http.StatusUnauthorized, errors.New("not found user")
	}
	return user, 0, nil
}
