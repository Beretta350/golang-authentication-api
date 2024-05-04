package controller

import (
	"net/http"

	"github.com/Beretta350/authentication/internal/app/service"
	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/Beretta350/authentication/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Login(c *gin.Context)
	Save(c *gin.Context)
}

type userController struct {
	service service.UserService
	jwt     jwt.JWTWrapper
}

func NewUserController(s service.UserService, j jwt.JWTWrapper) *userController {
	return &userController{service: s, jwt: j}
}

func (uc *userController) Login(c *gin.Context) {
	response := dto.OkResponse("User is already logged", nil)
	authHeader := c.GetHeader("Authorization")

	user := dto.UserRequest{}
	err := c.BindJSON(&user)
	if err != nil {
		c.Error(err)
		return
	}

	valid, _ := uc.jwt.ValidateToken(user.Username, authHeader)

	if !valid {
		response, err = uc.service.Login(c, user)
		if err != nil {
			c.Error(err)
			return
		}

		token, err := uc.jwt.GenerateJWT(user.Username)
		if err != nil {
			c.Error(err)
			return
		}

		c.Header("Authorization", token)
	}

	c.JSON(http.StatusOK, response)
}

func (uc *userController) Save(c *gin.Context) {
	user := dto.UserRequest{}
	err := c.BindJSON(&user)
	if err != nil {
		c.Error(err)
		return
	}

	response, err := uc.service.Save(c, user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
