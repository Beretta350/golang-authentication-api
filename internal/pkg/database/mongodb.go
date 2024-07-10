package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Beretta350/authentication/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	clientOpts *options.ClientOptions
}

func NewMongoDBClient() *MongoDB {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(getDbUri()).SetServerAPIOptions(serverAPI)
	return &MongoDB{clientOpts: opts}
}

func (db *MongoDB) ConnectDB(ctx context.Context) *mongo.Database {
	dbConfig := config.GetConfig().Database
	client, err := mongo.Connect(ctx, db.clientOpts)
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
