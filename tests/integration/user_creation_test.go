package integration_tests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Beretta350/authentication/internal/app/user/model"
	"github.com/Beretta350/authentication/tests"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_AuthenticationAPIUserCreation(t *testing.T) {
	ctx := context.Background()

	containers := tests.SetupContainers(t, ctx)
	defer containers.Terminate()

	createUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	createUserUsernameSubtests(t, ctx, containers.ApiPort, containers.MongoClient)
	createUserPasswordSubtests(t, ctx, containers.ApiPort, containers.MongoClient)
	createUserRoleSubtests(t, ctx, containers.ApiPort, containers.MongoClient)
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
