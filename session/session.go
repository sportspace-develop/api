package session

import (
	"fmt"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Session struct {
	sessions.Session
}

func New(c *gin.Context) Session {

	return Session{
		sessions.Default(c),
	}
}

// User
func (s *Session) GetUserId() uint {

	if s.Get("user_id") == nil {
		return 0
	}
	userId, err := strconv.Atoi(s.Get("user_id").(string))
	if err != nil {
		return 0
	}

	return uint(userId)
}

func (s *Session) SetUserId(id uint) {
	s.Set("user_id", fmt.Sprint(id))
}

//***************************
