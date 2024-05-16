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
	GetUserByID(c *gin.Context)
	Login(c *gin.Context)
	Save(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	RefreshToken(c *gin.Context)
}

const accessTokenName string = "accessToken"
const refreshTokenName string = "refreshToken"
const expireAccessTokenInSeconds int64 = 300    //5 min
const expireRefreshTokenInSeconds int64 = 86400 //24 hours

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
	c.JSON(http.StatusOK, dto.OkResponse("Success", userResponse))
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

	accessToken, err := jwt.GetJWTWrapper().GenerateJWT(user.ID, expireAccessTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	refreshToken, err := jwt.GetJWTWrapper().GenerateJWT(user.ID, expireRefreshTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(refreshTokenName, refreshToken, int(expireRefreshTokenInSeconds), "/refreshToken", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse("Login with success", gin.H{"accessToken": accessToken}))
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

	err = uc.service.Update(c, userReq)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Authorization", "")
	c.SetCookie(refreshTokenName, "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse("User successfully updated", nil))
}

func (uc *userController) Delete(c *gin.Context) {
	id := c.Query("id")

	err := uc.service.Delete(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Authorization", "")
	c.SetCookie(refreshTokenName, "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse("User successfully deleted", nil))
}

func (uc *userController) RefreshToken(c *gin.Context) {
	cookie, err := c.Request.Cookie(refreshTokenName)
	if err != nil || len(cookie.Value) <= 0 {
		c.Error(err)
		return
	}

	valid, userId, err := jwt.GetJWTWrapper().ValidateRefreshToken(cookie.Value)
	if !valid {
		c.Error(err)
		return
	}

	accessToken, err := jwt.GetJWTWrapper().GenerateJWT(userId, expireAccessTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	refreshToken, err := jwt.GetJWTWrapper().GenerateJWT(userId, expireRefreshTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(refreshTokenName, refreshToken, int(expireRefreshTokenInSeconds), "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse("Success", gin.H{"accessToken": accessToken}))
}
