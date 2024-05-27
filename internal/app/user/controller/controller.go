package controller

import (
	"net/http"

	messages_constants "github.com/Beretta350/authentication/internal/app/common/constants/messages"
	router_constants "github.com/Beretta350/authentication/internal/app/common/constants/router"
	token_constants "github.com/Beretta350/authentication/internal/app/common/constants/token"
	userModel "github.com/Beretta350/authentication/internal/app/user/model"
	userService "github.com/Beretta350/authentication/internal/app/user/service"
	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUserByID(c *gin.Context)
	Login(c *gin.Context)
	Save(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type userController struct {
	service userService.UserService
}

func NewUserController(s userService.UserService) *userController {
	return &userController{service: s}
}

func (uc *userController) GetUserByID(c *gin.Context) {
	userID := c.Query("id")

	user, err := uc.service.GetUserByID(c, userID)
	if err != nil {
		c.Error(err)
		return
	}

	userResponse := dto.NewUserResponseFromModel(*user)
	c.JSON(http.StatusOK, dto.OkResponse(messages_constants.SuccessMessage, userResponse))
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

	accessToken, refreshToken, err := uc.service.GenerateTokens(c, user.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(token_constants.RefreshTokenName, refreshToken, int(token_constants.ExpireRefreshTokenInSeconds), router_constants.RefreshTokenRoute, "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse(messages_constants.LoginSuccessMessage, gin.H{"accessToken": accessToken}))
}

func (uc *userController) Save(c *gin.Context) {
	userReq := userModel.User{}
	err := c.BindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = uc.service.Save(c, userReq)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.CreatedResponse("User successfully created", nil))
}

func (uc *userController) Update(c *gin.Context) {
	userReq := userModel.User{}
	err := c.BindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}
	userReq.ID = c.Query("id")

	err = uc.service.Update(c, userReq)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Authorization", "")
	c.SetCookie(token_constants.RefreshTokenName, "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse(messages_constants.UpdateSuccessMessage, nil))
}

func (uc *userController) Delete(c *gin.Context) {
	id := c.Query("id")

	err := uc.service.Delete(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Authorization", "")
	c.SetCookie(token_constants.RefreshTokenName, "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse(messages_constants.DeleteSuccessMessage, nil))
}

func (uc *userController) RefreshToken(c *gin.Context) {
	accessToken, refreshToken, err := uc.service.GenerateTokens(c, c.Query("id"))
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(token_constants.RefreshTokenName, refreshToken, int(token_constants.ExpireRefreshTokenInSeconds), "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse(messages_constants.RefreshTokenSuccessMessage, gin.H{"accessToken": accessToken}))
}
