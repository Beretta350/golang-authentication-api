package middleware

import (
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

		valid, _ := wrapper.ValidateAccessToken(userId, c.GetHeader("Authorization"))
		if !valid {
			c.Header("Authorization", "")
			defaultJWTErrorFunc(c, nil)
			return
		}

		c.Next()
	}
}

func defaultJWTErrorFunc(c *gin.Context, err error) {
	response := dto.UnauthorizedResponse("Invalid JWT token", err)
	c.JSON(response.StatusCode, response)
	c.Abort()
}
