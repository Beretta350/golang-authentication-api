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

func Test_AuthenticationAPIUserCreation(t *testing.T) {
	ctx := context.Background()

	mongoContainer, mongoIp, mongoPort := startMongoContainer(t, ctx)
	apiContainer, _, apiPort := startAPIContainer(t, ctx, mongoIp)
	mongoClient := newMongoClient(t, ctx, mongoPort)

	defer func() {
		mongoClient.Disconnect(ctx)
		apiContainer.Terminate(ctx)
		mongoContainer.Terminate(ctx)
	}()

	createUserHappyPathSubtest(t, ctx, apiPort, mongoClient)
	createUserUsernameSubtests(t, ctx, apiPort, mongoClient)
	createUserPasswordSubtests(t, ctx, apiPort, mongoClient)
	createUserRoleSubtests(t, ctx, apiPort, mongoClient)
}

func createUserHappyPathSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Creating new user happy path", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/save", apiPort)
		jsonStr := []byte(
			`{
				"username":"happypath",
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
		filter := bson.M{"username": "happypath"}
		err = mongoClient.Database("authentication").Collection("user").FindOne(ctx, filter).Decode(&user)
		assert.NoError(t, err)

		assert.Equal(t, user.Username, "happypath")
		assert.Equal(t, user.Roles, []string{"USER"})
	})
}

func createUserUsernameSubtests(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Create new user with no username", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/save", apiPort)
		jsonStr := []byte(
			`{
				"password": "ABCD12345",
				"roles": ["USER"]
			}`,
		)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		err = resp.Body.Close()
		assert.NoError(t, err)

		user := model.User{}
		filter := bson.M{"username": ""}
		err = mongoClient.Database("authentication").Collection("user").FindOne(ctx, filter).Decode(&user)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
	})
}

func createUserPasswordSubtests(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Create new user with no password", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/save", apiPort)
		jsonStr := []byte(
			`{
				"username":"wrongpassword",
				"roles": ["USER"]
			}`,
		)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		err = resp.Body.Close()
		assert.NoError(t, err)

		user := model.User{}
		filter := bson.M{"username": "wrongpassword"}
		err = mongoClient.Database("authentication").Collection("user").FindOne(ctx, filter).Decode(&user)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
	})
	t.Run("Create new user with less than 8 characters in the password", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/save", apiPort)
		jsonStr := []byte(
			`{
				"username":"wrongpassword",
				"password": "ABC",
				"roles": ["USER"]
			}`,
		)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		err = resp.Body.Close()
		assert.NoError(t, err)

		user := model.User{}
		filter := bson.M{"username": "wrongpassword"}
		err = mongoClient.Database("authentication").Collection("user").FindOne(ctx, filter).Decode(&user)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
	})
}

func createUserRoleSubtests(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Create new user with wrong roles", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/save", apiPort)
		jsonStr := []byte(
			`{
				"username":"wrongrole",
				"password": "ABCDE123456",
				"roles": ["WROONG"]
			}`,
		)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		err = resp.Body.Close()
		assert.NoError(t, err)

		user := model.User{}
		filter := bson.M{"username": "wrongrole"}
		err = mongoClient.Database("authentication").Collection("user").FindOne(ctx, filter).Decode(&user)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
	})
	t.Run("Create new user with no roles", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/save", apiPort)
		jsonStr := []byte(
			`{
				"username":"wrongrole",
				"password": "ABCDE123456"
			}`,
		)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		err = resp.Body.Close()
		assert.NoError(t, err)

		user := model.User{}
		filter := bson.M{"username": "wrongrole"}
		err = mongoClient.Database("authentication").Collection("user").FindOne(ctx, filter).Decode(&user)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
	})
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
			Context: "..",
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
