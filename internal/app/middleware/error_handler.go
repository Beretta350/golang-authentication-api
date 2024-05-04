package middleware

import (
	"errors"
	"strings"

	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		var validationErrors validator.ValidationErrors
		var response *dto.ResponseMessage

		err := c.Errors.Last().Err
		switch {
		case errors.As(err, &validationErrors):
			var errs []string
			for _, e := range validationErrors {
				errMsg := strings.Split(e.Error(), "Error:")[1]
				errs = append(errs, errMsg)
			}
			response = dto.BadRequestResponse("Invalid data", errs)
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			response = dto.UnauthorizedResponse("Invalid username or password")
		default:
			response = dto.InternalErrorResponse(err.Error(), nil)
		}
		c.JSON(response.StatusCode, response)
	}
}
