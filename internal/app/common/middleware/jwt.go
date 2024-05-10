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

		var requestBody struct {
			Username string `json:"username"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			defaultJWTErrorFunc(c, err)
			return
		}

		valid, err := wrapper.ValidateAccessToken(requestBody.Username, c.GetHeader("Authorization"))
		if !valid {
			c.Header("Authorization", "")
			defaultJWTErrorFunc(c, err)
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
