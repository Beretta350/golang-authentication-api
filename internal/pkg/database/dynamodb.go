package database

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDB struct {
	session *session.Session
}

func NewDynamoDBSession() *DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Fatalln("failed to create DynamoDB session: %w", err)
	}

	return &DynamoDB{session: sess}
}

func (db *DynamoDB) ConnectDB(ctx context.Context) *dynamodb.DynamoDB {
	svc := dynamodb.New(db.session)
	input := &dynamodb.ListTablesInput{}

	//List tables
	_, err := svc.ListTables(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Fatalln("Error while connecting to DynamoDB: %w", aerr)
		}
		log.Fatalln("Error while connecting to DynamoDB: %w", err)
	}

	return svc
}

func (db *DynamoDB) GetURI() string {
	return getDbUri()
}
