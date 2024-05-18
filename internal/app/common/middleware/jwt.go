package middleware

import (
	"github.com/Beretta350/authentication/internal/app/common/enum/constants"
	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/Beretta350/authentication/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func JWTHandler(wrapper jwt.JWTWrapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		if wrapper.IsIgnoredPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		userId := c.Query("id")
		if len(userId) <= 0 {
			defaultJWTErrorFunc(c, nil)
			return
		}

		valid, _ := wrapper.ValidateToken(userId, c.GetHeader("Authorization"))
		if !valid {
			c.Header("Authorization", "")
			defaultJWTErrorFunc(c, nil)
			return
		}

		if c.Request.URL.Path == constants.RefreshTokenRoute {
			validateRefreshTokenCookie(c, wrapper, userId)
		}

		c.Next()
	}
}

func validateRefreshTokenCookie(c *gin.Context, wrapper jwt.JWTWrapper, userId string) {
	cookie, err := c.Request.Cookie(constants.RefreshTokenName)
	if err != nil || len(cookie.Value) <= 0 {
		c.Header("Authorization", "")
		defaultJWTErrorFunc(c, nil)
		return
	}

	valid, _ := wrapper.ValidateToken(userId, cookie.Value)
	if !valid {
		c.Header("Authorization", "")
		defaultJWTErrorFunc(c, nil)
		return
	}
}

func defaultJWTErrorFunc(c *gin.Context, err error) {
	response := dto.UnauthorizedResponse("Invalid JWT token", err)
	c.JSON(response.StatusCode, response)
	c.Abort()
}
