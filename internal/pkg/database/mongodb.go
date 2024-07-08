package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Beretta350/authentication/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct{}

func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

func (db *MongoDB) ConnectDB(ctx context.Context) *mongo.Database {
	dbConfig := config.GetConfig().Database
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(db.GetURI()).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("error pinging MongoDB server: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client.Database(dbConfig.GetDatabase())
}

func (db *MongoDB) GetURI() string {
	return getDbUri()
}
