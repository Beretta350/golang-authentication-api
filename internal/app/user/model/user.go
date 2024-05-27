package model

import (
	"github.com/Beretta350/authentication/internal/app/common/constants/roles"
	commonModel "github.com/Beretta350/authentication/internal/app/common/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	commonModel.Base `bson:",inline"`
	Username         string   `json:"username" bson:"username" validate:"required,min=3"`
	Password         string   `json:"password" bson:"password" validate:"required,min=8"`
	Roles            []string `json:"roles,omitempty" bson:"roles" validate:"required,gt=0,validRoles"`
}

func NewUserModel(username, password string, roles []string) *User {
	id := uuid.NewString()
	return &User{Base: commonModel.Base{ID: id}, Username: username, Password: password, Roles: roles}
}

func (u User) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("validRoles", areValidRoles)
	return validate.Struct(u)
}

func areValidRoles(fl validator.FieldLevel) bool {
	// name := fl.FieldName()
	userRoles := fl.Field().Interface().([]string)

	validRoles := map[string]bool{
		roles.ADMIN: true, roles.USER: true,
	}

	for _, role := range userRoles {
		if !validRoles[role] {
			return false
		}
	}
	return true
}
