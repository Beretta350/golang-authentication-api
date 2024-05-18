package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Beretta350/authentication/internal/app/common/enum/constants"
	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/Beretta350/authentication/tests"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestAuthenticationAPIUser_Refresh(t *testing.T) {
	ctx := context.Background()

	containers := tests.SetupContainers(t, ctx)
	defer containers.Terminate()

	//Create the user to make login
	createUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	//Login to register the accessToken
	loginUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	refreshUserHappyPathSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	refreshNoHeaderSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
	refreshNoCookieSubtest(t, ctx, containers.ApiPort, containers.MongoClient)
}

func refreshUserHappyPathSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Refresh user token happy path", func(t *testing.T) {
		userId := tests.GetUserIdFromToken(t, accessToken)
		url := fmt.Sprintf("http://localhost:%s/refreshToken?id=%s", apiPort, userId)

		req, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", accessToken)
		req.AddCookie(&http.Cookie{Name: constants.RefreshTokenName, Value: refreshToken})

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		body := dto.ResponseMessage{}
		err = json.Unmarshal(bodyBytes, &body)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, "Token refreshed Successfully", body.Message)
		err = resp.Body.Close()
		assert.NoError(t, err)

		val := body.Data.(map[string]interface{})
		accessToken = val["accessToken"].(string)
		refreshToken = resp.Cookies()[0].Value
	})
}

func refreshNoHeaderSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Refresh user token no authentication header", func(t *testing.T) {
		userId := tests.GetUserIdFromToken(t, accessToken)
		url := fmt.Sprintf("http://localhost:%s/refreshToken?id=%s", apiPort, userId)

		req, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: constants.RefreshTokenName, Value: refreshToken})

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

func refreshNoCookieSubtest(t *testing.T, ctx context.Context, apiPort string, mongoClient *mongo.Client) {
	t.Run("Refresh user token no refresh coockie token", func(t *testing.T) {
		userId := tests.GetUserIdFromToken(t, accessToken)
		url := fmt.Sprintf("http://localhost:%s/refreshToken?id=%s", apiPort, userId)

		req, err := http.NewRequest("GET", url, nil)
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
