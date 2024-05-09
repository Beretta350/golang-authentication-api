package controller

import (
	"net/http"

	userModel "github.com/Beretta350/authentication/internal/app/user/model"
	userService "github.com/Beretta350/authentication/internal/app/user/service"
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
	service userService.UserService
	jwt     jwt.JWTWrapper
}

func NewUserController(s userService.UserService) *userController {
	return &userController{service: s}
}

func (uc *userController) Login(c *gin.Context) {
	userReq := userModel.User{}
	err := c.BindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := uc.service.Login(c, userReq)
	if err != nil {
		c.Error(err)
		return
	}

	token, err := jwt.GetJWTWrapper().GenerateJWT(user.Username)
	if err != nil {
		c.Error(err)
		return
	}

	userResponse := dto.NewUserResponseFromModel(*user)

	c.Header("Authorization", token)
	c.JSON(http.StatusOK, dto.OkResponse("Login with success", userResponse))
}

func (uc *userController) Save(c *gin.Context) {
	userReq := userModel.User{}
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

	userReq := userModel.User{}
	err := c.BindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}
	userReq.ID = c.Query("id")

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
	username := c.Query("username")

	err := uc.service.Delete(c, username)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.OkResponse("User successfully deleted", nil))
}
