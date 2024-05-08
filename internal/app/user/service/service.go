package service

import (
	"context"

	"github.com/Beretta350/authentication/internal/app/user/model"
	userRepo "github.com/Beretta350/authentication/internal/app/user/repository"
	"github.com/Beretta350/authentication/internal/pkg/crypto"
)

type UserService interface {
	Login(ctx context.Context, userReq model.User) (*model.User, error)
	Save(ctx context.Context, userReq model.User) (*model.User, error)
	Update(ctx context.Context, userReq model.User) (*model.User, error)
	Delete(ctx context.Context, username string) error
}

type userService struct {
	repo userRepo.UserRepository
}

func NewUserService(r userRepo.UserRepository) *userService {
	return &userService{repo: r}
}

func (us *userService) Login(ctx context.Context, userReq model.User) (*model.User, error) {

	//Find
	user, err := us.repo.FindByUsername(ctx, userReq.Username)
	if err != nil {
		return nil, err
	}

	//Check password
	err = crypto.CheckPassword(userReq.Password, []byte(user.Password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) Save(ctx context.Context, userReq model.User) (*model.User, error) {
	user := model.NewUserModel(userReq.Username, userReq.Password, userReq.Roles)
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

func (us *userService) Update(ctx context.Context, userReq model.User) (*model.User, error) {
	user, err := us.repo.FindByUsername(ctx, userReq.Username)
	if err != nil {
		return nil, err
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	encryptedData, err := crypto.EncryptPassword(userReq.Password)
	if err != nil {
		return nil, err
	}

	user.Password = string(encryptedData)

	updatedUser, err := us.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (us *userService) Delete(ctx context.Context, username string) error {
	user, err := us.repo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	err = user.Validate()
	if err != nil {
		return err
	}

	err = us.repo.Delete(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
