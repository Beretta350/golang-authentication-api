package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Beretta350/authentication/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context, cfg config.DatabaseConfig) *mongo.Database {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(cfg.GetURI()).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("error pinging MongoDB server: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client.Database(cfg.Database)
}
