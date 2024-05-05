package service

import (
	"context"

	"github.com/Beretta350/authentication/internal/app/model"
	"github.com/Beretta350/authentication/internal/app/repository"
	"github.com/Beretta350/authentication/internal/pkg/crypto"
	"github.com/Beretta350/authentication/internal/pkg/dto"
)

type UserService interface {
	Save(ctx context.Context, request dto.UserRequest) (*model.User, error)
	Login(ctx context.Context, request dto.UserRequest) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *userService {
	return &userService{repo: r}
}

func (us *userService) Login(ctx context.Context, request dto.UserRequest) (*model.User, error) {

	//Find
	user, err := us.repo.FindByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	//Check password
	err = crypto.CheckPassword(request.Password, []byte(user.Password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) Save(ctx context.Context, request dto.UserRequest) (*model.User, error) {
	user := model.NewUserModel(request)
	err := user.Validate()
	if err != nil {
		return nil, err
	}

	encryptedData, err := crypto.EncryptPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = string(encryptedData)

	err = us.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
