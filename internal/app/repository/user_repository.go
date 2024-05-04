package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Beretta350/authentication/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Save(ctx context.Context, user *model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(d *mongo.Database) *userRepository {
	return &userRepository{collection: d.Collection("user")}
}

func (ur *userRepository) Save(ctx context.Context, user *model.User) error {

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	filter := bson.M{"username": username}
	err := ur.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
