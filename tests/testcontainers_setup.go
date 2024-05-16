package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestContainersInfos struct {
	MongoClient *mongo.Client
	ApiPort     string
	Terminate   func()
}

func SetupContainers(t *testing.T, ctx context.Context) *TestContainersInfos {
	mongoContainer, mongoIp, mongoPort := startMongoContainer(t, ctx)
	apiContainer, _, apiPort := startAPIContainer(t, ctx, mongoIp)
	mongoClient := newMongoClient(t, ctx, mongoPort)

	return &TestContainersInfos{
		MongoClient: mongoClient,
		ApiPort:     apiPort,
		Terminate: func() {
			mongoClient.Disconnect(ctx)
			apiContainer.Terminate(ctx)
			mongoContainer.Terminate(ctx)
		},
	}
}

func startAPIContainer(t *testing.T, ctx context.Context, mongoHost string) (testcontainers.Container, string, string) {
	req := testcontainers.ContainerRequest{
		Name: "authentication",
		Env: map[string]string{
			"SERVER_PORT": "8080",
			"SERVER_MODE": "debug",
			"DB_DRIVER":   "mongodb",
			"DB_HOST":     mongoHost,
			"DB_PORT":     "27017",
			"DB_USERNAME": "root",
			"DB_PASSWORD": "root",
			"DB_DATABASE": "authentication",
			"JWT_SECRET":  "Hxj1pW48QqcnSQAc5",
		},
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "../../",
		},
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForListeningPort("8080/tcp"),
	}

	apiContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	apiContainer.Start(ctx)
	assert.NoError(t, err)
	apiIp, err := apiContainer.ContainerIP(ctx)
	assert.NoError(t, err)
	apiPort, err := apiContainer.MappedPort(ctx, "8080")
	assert.NoError(t, err)

	return apiContainer, apiIp, apiPort.Port()
}

func startMongoContainer(t *testing.T, ctx context.Context) (testcontainers.Container, string, string) {
	req := testcontainers.ContainerRequest{
		Image: "mongo:latest",
		Name:  "mongodb",
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "root",
			"MONGO_INITDB_DATABASE":      "authentication",
		},
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("MongoDB init process complete; ready for start up."),
	}

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	mongoContainer.Start(ctx)
	assert.NoError(t, err)
	mongoIp, err := mongoContainer.ContainerIP(ctx)
	assert.NoError(t, err)
	mongoPort, err := mongoContainer.MappedPort(ctx, "27017")
	assert.NoError(t, err)

	return mongoContainer, mongoIp, mongoPort.Port()
}

func newMongoClient(t *testing.T, ctx context.Context, mongoContainerPort string) *mongo.Client {
	mongoURI := fmt.Sprintf("mongodb://root:root@localhost:%s", mongoContainerPort)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	assert.NoError(t, err)

	err = mongoClient.Ping(ctx, nil)
	assert.NoError(t, err)
	return mongoClient
}
