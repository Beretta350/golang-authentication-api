package repository

import (
	"context"
	"time"

	userModel "github.com/Beretta350/authentication/internal/app/user/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type userDynamoDBRepository struct {
	database  *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBUserRepository(d *dynamodb.DynamoDB) *userDynamoDBRepository {
	return &userDynamoDBRepository{
		database:  d,
		tableName: "user",
	}
}

func (ur *userDynamoDBRepository) Save(ctx context.Context, user *userModel.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = ur.database.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(ur.tableName),
		Item:      item,
	})
	return err
}

func (ur *userDynamoDBRepository) Update(ctx context.Context, user *userModel.User) error {
	user.UpdatedAt = time.Now()

	//update expression attributes
	update := map[string]*dynamodb.AttributeValue{
		":username": {S: aws.String(user.Username)},
		":updateat": {S: aws.String(user.UpdatedAt.Format(time.RFC3339))},
	}

	_, err := ur.database.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(ur.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(user.ID)},
		},
		ExpressionAttributeValues: update,
		UpdateExpression:          aws.String("SET username = :username, updatedat = :updatedat"),
		ReturnValues:              aws.String("UPDATED_NEW"),
	})
	return err
}

func (ur *userDynamoDBRepository) Delete(ctx context.Context, user *userModel.User) error {
	_, err := ur.database.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(ur.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(user.ID)},
		},
	})
	return err
}

func (ur *userDynamoDBRepository) FindByUsername(ctx context.Context, username string) (*userModel.User, error) {
	var user userModel.User
	result, err := ur.database.Query(&dynamodb.QueryInput{
		TableName:              aws.String(ur.tableName),
		IndexName:              aws.String("Username-index"), // Assume there is a GSI on Username
		KeyConditionExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {S: aws.String(username)},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	return &user, err
}

func (ur *userDynamoDBRepository) FindByID(ctx context.Context, id string) (*userModel.User, error) {
	var user userModel.User
	result, err := ur.database.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ur.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	return &user, err
}
