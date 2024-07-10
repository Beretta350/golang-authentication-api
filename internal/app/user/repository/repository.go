package repository

import (
	"context"

	userModel "github.com/Beretta350/authentication/internal/app/user/model"
)

type UserRepository interface {
	Save(ctx context.Context, user *userModel.User) error
	Update(ctx context.Context, user *userModel.User) error
	Delete(ctx context.Context, user *userModel.User) error
	FindByUsername(ctx context.Context, username string) (*userModel.User, error)
	FindByID(ctx context.Context, id string) (*userModel.User, error)
}
