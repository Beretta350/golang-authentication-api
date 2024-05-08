package dto

import "github.com/Beretta350/authentication/internal/app/user/model"

type UserResponse struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles,omitempty"`
}

func NewUserResponseFromModel(user model.User) *UserResponse {
	return &UserResponse{ID: user.ID, Username: user.Username, Roles: user.Roles}
}
