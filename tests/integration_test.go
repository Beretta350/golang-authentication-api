package tests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Beretta350/authentication/internal/app/user/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuthenticationAPI(t *testing.T) {
	ctx := context.Background()

	mongoContainer, mongoIp, mongoPort := startMongoContainer(t, ctx)
	apiContainer, _, apiPort := startAPIContainer(t, ctx, mongoIp)
	mongoClient := newMongoClient(t, ctx, mongoPort)

	t.Run("Creating new user happy path", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/save", apiPort)
		jsonStr := []byte(
			`{
				"username":"someone",
				"password": "ABCD1234",
				"roles": ["USER"]
			}`,
		)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusCreated)
		err = resp.Body.Close()
		assert.NoError(t, err)

		user := model.User{}
		filter := bson.M{"username": "someone"}
		err = mongoClient.Database("authentication").Collection("user").FindOne(ctx, filter).Decode(&user)
		assert.NoError(t, err)

		assert.Equal(t, user.Username, "someone")
		assert.Equal(t, user.Roles, []string{"USER"})
	})

	err := mongoClient.Disconnect(ctx)
	assert.NoError(t, err)

	err = apiContainer.Terminate(ctx)
	assert.NoError(t, err)

	err = mongoContainer.Terminate(ctx)
	assert.NoError(t, err)
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
		Image:        "authentication-api",
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
