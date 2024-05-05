package middleware

import (
	"errors"
	"strings"

	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		var response *dto.ResponseMessage
		var validationErrors validator.ValidationErrors

		err := c.Errors.Last().Err
		switch {
		case errors.As(err, &validationErrors):
			response = multipleErrorMessages(validationErrors)
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			response = dto.UnauthorizedResponse("Invalid username or password", nil)
		default:
			response = dto.InternalErrorResponse(err.Error(), nil)

			if _, ok := err.(*jwt.ValidationError); ok {
				response = dto.UnauthorizedResponse("Invalid token", err.Error())
			}
		}
		c.JSON(response.StatusCode, response)
	}
}

func multipleErrorMessages(validationErrors validator.ValidationErrors) *dto.ResponseMessage {
	var errs []string
	for _, e := range validationErrors {
		errMsg := strings.Split(e.Error(), "Error:")[1]
		errs = append(errs, errMsg)
	}
	return dto.BadRequestResponse("Invalid data", errs)
}
