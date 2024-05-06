package controller

import (
	"errors"
	"net/http"

	"github.com/Beretta350/authentication/internal/app/model"
	"github.com/Beretta350/authentication/internal/app/service"
	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/Beretta350/authentication/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Login(c *gin.Context)
	Save(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
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

	userReq := model.User{}
	err := c.BindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	if len(authHeader) > 0 {
		loged := uc.loginWithToken(c, userReq.Username, authHeader)
		if loged {
			c.JSON(http.StatusOK, dto.OkResponse("User is already logged", nil))
		}
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

	userResponse := dto.NewUserResponseFromModel(*user)

	c.Header("Authorization", token)
	c.JSON(http.StatusOK, dto.OkResponse("Login with success", userResponse))
}

func (uc *userController) Save(c *gin.Context) {
	userReq := model.User{}
	err := c.BindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	_, err = uc.service.Save(c, userReq)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.CreatedResponse("User successfully created", nil))
}

func (uc *userController) Update(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	userReq := model.User{}
	err := c.BindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}
	userReq.ID = c.Query("id")

	loged := uc.loginWithToken(c, userReq.Username, authHeader)
	if !loged {
		return
	}

	user, err := uc.service.Update(c, userReq)
	if err != nil {
		c.Error(err)
		return
	}

	token, err := uc.jwt.GenerateJWT(user.Username)
	if err != nil {
		c.Error(err)
		return
	}

	userResponse := dto.NewUserResponseFromModel(*user)
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, dto.OkResponse("User successfully updated", userResponse))
}

func (uc *userController) Delete(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	username := c.Query("username")

	loged := uc.loginWithToken(c, username, authHeader)
	if !loged {
		return
	}

	err := uc.service.Delete(c, username)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.OkResponse("User successfully deleted", nil))
}

func (uc *userController) loginWithToken(c *gin.Context, username, authHeader string) bool {
	valid, err := uc.jwt.ValidateToken(username, authHeader)
	if !valid {
		c.Header("Authorization", "")

		if err != nil {
			c.Error(err)
		} else {
			c.Error(errors.New("invalid token"))
		}

		return false
	}

	return true
}
