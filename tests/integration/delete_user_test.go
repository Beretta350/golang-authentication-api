package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Beretta350/authentication/internal/app/user/model"
	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/Beretta350/authentication/tests"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestAuthenticationAPIUser_Delete(t *testing.T) {
	ctx := context.Background()

	containers := tests.SetupContainers(t, ctx)
	defer containers.Terminate()

	//Create the user to make login
	createUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	//Login to register the accessToken
	loginUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	deleteUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	deleteNoHeaderSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	deleteWrongIdSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
}

func deleteUserHappyPathSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Delete user happy path", func(t *testing.T) {

		userId := tests.GetUserIdFromToken(t, accessToken)

		url := fmt.Sprintf("http://localhost:%s/delete?id=%s", apiPort, userId)

		req, err := http.NewRequest("DELETE", url, nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		body := dto.ResponseMessage{}
		err = json.Unmarshal(bodyBytes, &body)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, "User successfully deleted", body.Message)

		err = resp.Body.Close()
		assert.NoError(t, err)

		user := model.User{}
		filter := bson.M{"username": "happypath"}
		err = mongoClient.Database("authentication").Collection("user").FindOne(ctx, filter).Decode(&user)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
	})
}

func deleteNoHeaderSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Delete user no authorization header", func(t *testing.T) {

		userId := tests.GetUserIdFromToken(t, accessToken)

		url := fmt.Sprintf("http://localhost:%s/delete?id=%s", apiPort, userId)

		req, err := http.NewRequest("DELETE", url, nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		body := dto.ResponseMessage{}
		err = json.Unmarshal(bodyBytes, &body)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
		assert.Equal(t, "Invalid JWT token", body.Message)

		err = resp.Body.Close()
		assert.NoError(t, err)
	})
}

func deleteWrongIdSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Delete user wrong ID", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/delete?id=%s", apiPort, "123456789")

		req, err := http.NewRequest("DELETE", url, nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		body := dto.ResponseMessage{}
		err = json.Unmarshal(bodyBytes, &body)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
		assert.Equal(t, "Invalid JWT token", body.Message)

		err = resp.Body.Close()
		assert.NoError(t, err)
	})
}
