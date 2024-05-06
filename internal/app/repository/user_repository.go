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
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, user *model.User) error
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

func (ur *userRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	filter := bson.M{"_id": user.ID, "username": user.Username}
	update := bson.M{"$set": user}

	user.UpdatedAt = time.Now()

	result, err := ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount <= 0 {
		return nil, errors.New("no updated users")
	}

	return user, nil
}

func (ur *userRepository) Delete(ctx context.Context, user *model.User) error {
	filter := bson.M{"_id": user.ID, "username": user.Username}
	_, err := ur.collection.DeleteOne(ctx, filter)
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
