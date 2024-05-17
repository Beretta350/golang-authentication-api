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

var accessToken string
var refreshToken string

func TestAuthenticationAPIUser_Login(t *testing.T) {
	ctx := context.Background()

	containers := tests.SetupContainers(t, ctx)
	defer containers.Terminate()

	//Create the user to make login
	createUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	loginUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	loginUserUsernameSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	loginUserPasswordSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
}

func loginUserHappyPathSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Login user happy path", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/login", apiPort)
		jsonStr := []byte(
			`{
				"username":"happypath",
				"password": "ABCD1234"
			}`,
		)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, "Login with success", body.Message)
		err = resp.Body.Close()
		assert.NoError(t, err)

		val := body.Data.(map[string]interface{})
		accessToken = val["accessToken"].(string)
		refreshToken = resp.Cookies()[0].Value
	})
}

func loginUserUsernameSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Login user wrong username", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/login", apiPort)
		jsonStr := []byte(
			`{
				"username":"happypath",
				"password": "12345ABCD"
			}`,
		)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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
		assert.Equal(t, "Invalid username or password", body.Message)
		err = resp.Body.Close()
		assert.NoError(t, err)
	})
}

func loginUserPasswordSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Login user wrong password", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:%s/login", apiPort)
		jsonStr := []byte(
			`{
				"username":"wroong",
				"password": "ABCD1234"
			}`,
		)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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
		assert.Equal(t, "invalid username or password", body.Message)
		err = resp.Body.Close()
		assert.NoError(t, err)
	})
}
