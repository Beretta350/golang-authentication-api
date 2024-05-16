package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/Beretta350/authentication/tests"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestAuthenticationAPIUser_Update(t *testing.T) {
	ctx := context.Background()

	containers := tests.SetupContainers(t, ctx)
	defer containers.Terminate()

	//Create the user to make login
	createUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	//Login to register the accessToken
	loginUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	updateUserMissingDataSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	updateUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	updateNoHeaderSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	updateWrongIdSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
}

func updateUserHappyPathSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Update user happy path", func(t *testing.T) {

		userId := tests.GetUserIdFromToken(t, accessToken)

		url := fmt.Sprintf("http://localhost:%s/update?id=%s", apiPort, userId)
		jsonStr := []byte(
			`{
				"username":"updateSuccess",
				"password": "12345ABCD"
			}`,
		)

		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
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
		assert.Equal(t, "User successfully updated", body.Message)

		err = resp.Body.Close()
		assert.NoError(t, err)
	})
}

func updateUserMissingDataSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Update user missing data", func(t *testing.T) {

		userId := tests.GetUserIdFromToken(t, accessToken)

		url := fmt.Sprintf("http://localhost:%s/update?id=%s", apiPort, userId)
		jsonStr := []byte(
			`{}`,
		)

		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
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

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, "missing data in request", body.Message)

		err = resp.Body.Close()
		assert.NoError(t, err)
	})
}

func updateNoHeaderSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Update user no authorization header", func(t *testing.T) {

		userId := tests.GetUserIdFromToken(t, accessToken)

		url := fmt.Sprintf("http://localhost:%s/update?id=%s", apiPort, userId)
		jsonStr := []byte(
			`{
				"username":"updateSuccess",
				"password": "12345ABCD"
			}`,
		)

		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
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

func updateWrongIdSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Update user wrong ID", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/update?id=%s", apiPort, "123456789")
		jsonStr := []byte(
			`{
				"username":"updateSuccess",
				"password": "12345ABCD"
			}`,
		)

		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
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
