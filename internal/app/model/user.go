package model

import (
	"github.com/Beretta350/authentication/internal/app/enum"
	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	Base     `bson:",inline"`
	Username string   `bson:"username" validate:"required,min=3"`
	Password string   `bson:"password" validate:"required,min=8"`
	Roles    []string `bson:"roles" validate:"required,gt=0,validRoles"`
}

func NewUserModel(request dto.UserRequest) *User {
	id := uuid.NewString()
	return &User{Base: Base{ID: id}, Username: request.Username, Password: request.Password, Roles: request.Roles}
}

func (u User) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("validRoles", areValidRoles)
	return validate.Struct(u)
}

func areValidRoles(fl validator.FieldLevel) bool {
	// name := fl.FieldName()
	roles := fl.Field().Interface().([]string)

	validRoles := map[string]bool{
		enum.ADMIN: true, enum.USER: true,
	}

	for _, role := range roles {
		if !validRoles[role] {
			return false
		}
	}
	return true
}
