package controller

import (
	"github.com/Beretta350/authentication/pkg/csrf"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CSRFController interface {
	GetToken(c *gin.Context)
}

type csrfController struct{}

func NewCSRFController() *csrfController {
	return &csrfController{}
}

func (cs *csrfController) GetToken(c *gin.Context) {
	session := sessions.Default(c)
	secret := c.MustGet(csrf.Secret).(string)

	if _, ok := c.Get(csrf.Token); ok {
		return
	}

	s, ok := session.Get(csrf.Session).(string)
	if !ok {
		s = uuid.New().String()
		session.Set(csrf.Session, s)
		session.Save()
	}
	token := csrf.GetCSRFWrapper().GenerateToken(secret, s)
	c.Set(csrf.Token, token)
}
