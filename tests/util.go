package tests

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetUserIdFromToken(t *testing.T, token string) string {
	tokenParts := strings.Split(token, ".")

	payload, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	assert.NoError(t, err)

	var claims map[string]interface{}
	err = json.Unmarshal(payload, &claims)
	assert.NoError(t, err)

	userID, ok := claims["id"].(string)
	assert.True(t, ok)

	return userID
}
