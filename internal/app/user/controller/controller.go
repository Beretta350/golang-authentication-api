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
	GetUser(c *gin.Context)
	Login(c *gin.Context)
	Save(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	RefreshToken(c *gin.Context)
}

const refreshTokenName string = "refresh_token"
const accessTokenName string = "access_token"
const expireAccessTokenInSeconds int64 = 1800   //30 min
const expireRefreshTokenInSeconds int64 = 86400 //24 hours

type userController struct {
	service userService.UserService
	jwt     jwt.JWTWrapper
}

func NewUserController(s userService.UserService) *userController {
	return &userController{service: s}
}

func (uc *userController) GetUser(c *gin.Context) {

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

	accessToken, err := jwt.GetJWTWrapper().GenerateJWT(user.Username, expireAccessTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	refreshToken, err := jwt.GetJWTWrapper().GenerateJWT(user.Username, expireRefreshTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Authorization", accessToken)
	c.SetCookie(refreshTokenName, refreshToken, int(expireRefreshTokenInSeconds), "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse("Login with success", gin.H{"access_token": accessToken}))
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

	accessToken, err := jwt.GetJWTWrapper().GenerateJWT(user.Username, expireAccessTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	refreshToken, err := jwt.GetJWTWrapper().GenerateJWT(user.Username, expireRefreshTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	userResponse := dto.NewUserResponseFromModel(*user)

	c.Header("Authorization", accessToken)
	c.SetCookie(refreshTokenName, refreshToken, int(expireRefreshTokenInSeconds), "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse("User successfully updated", userResponse))
}

func (uc *userController) Delete(c *gin.Context) {
	username := c.Query("username")

	err := uc.service.Delete(c, username)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Authorization", "")
	c.SetCookie(refreshTokenName, "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse("User successfully deleted", nil))
}

func (uc *userController) RefreshToken(c *gin.Context) {
	cookie, err := c.Request.Cookie(accessTokenName)
	if err != nil || len(cookie.Value) <= 0 {
		c.Error(err)
		return
	}

	valid, username, err := jwt.GetJWTWrapper().ValidateRefreshToken(cookie.Value)
	if !valid {
		c.Error(err)
		return
	}

	accessToken, err := jwt.GetJWTWrapper().GenerateJWT(username, expireAccessTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	refreshToken, err := jwt.GetJWTWrapper().GenerateJWT(username, expireRefreshTokenInSeconds)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(refreshTokenName, refreshToken, int(expireRefreshTokenInSeconds), "/", "localhost", false, true)
	c.JSON(http.StatusOK, dto.OkResponse("Login with success", gin.H{"access_token": accessToken}))
}
