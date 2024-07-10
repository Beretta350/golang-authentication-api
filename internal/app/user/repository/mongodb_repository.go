package repository

import (
	"context"
	"errors"
	"time"

	userModel "github.com/Beretta350/authentication/internal/app/user/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoDBRepository struct {
	collection *mongo.Collection
}

func NewMongoDBUserRepository(d *mongo.Database) *userMongoDBRepository {
	return &userMongoDBRepository{collection: d.Collection("user")}
}

func (ur *userMongoDBRepository) Save(ctx context.Context, user *userModel.User) error {

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userMongoDBRepository) Update(ctx context.Context, user *userModel.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	user.UpdatedAt = time.Now()

	result, err := ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount <= 0 {
		return errors.New("no updated users")
	}

	return nil
}

func (ur *userMongoDBRepository) Delete(ctx context.Context, user *userModel.User) error {
	filter := bson.M{"_id": user.ID}
	_, err := ur.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userMongoDBRepository) FindByUsername(ctx context.Context, username string) (*userModel.User, error) {
	var user userModel.User
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

func (ur *userMongoDBRepository) FindByID(ctx context.Context, id string) (*userModel.User, error) {
	var user userModel.User
	filter := bson.M{"_id": id}
	err := ur.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
