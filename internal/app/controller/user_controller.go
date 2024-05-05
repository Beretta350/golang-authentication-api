package controller

import (
	"errors"
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
	authHeader := c.GetHeader("Authorization")

	userReq := dto.UserRequest{}
	err := c.BindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	if len(authHeader) > 0 {
		uc.loginWithToken(c, userReq.Username, authHeader)
		return
	}

	user, err := uc.service.Login(c, userReq)
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
	c.JSON(http.StatusOK, dto.OkResponse("Login with success", user))
}

func (uc *userController) Save(c *gin.Context) {
	user := dto.UserRequest{}
	err := c.BindJSON(&user)
	if err != nil {
		c.Error(err)
		return
	}

	_, err = uc.service.Save(c, user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.CreatedResponse("User created with success", nil))
}

func (uc *userController) loginWithToken(c *gin.Context, username, authHeader string) {
	valid, err := uc.jwt.ValidateToken(username, authHeader)
	if !valid {
		c.Header("Authorization", "")

		if err != nil {
			c.Error(err)
		} else {
			c.Error(errors.New("invalid token"))
		}

		return
	}

	c.JSON(http.StatusOK, dto.OkResponse("User is already logged", nil))
}
