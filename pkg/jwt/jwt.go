package jwt

import (
	"sync"
	"time"

	"github.com/Beretta350/authentication/pkg/util"
	"github.com/golang-jwt/jwt"
)

type JWTWrapper interface {
	GenerateJWT(username string, expire int64) (string, error)
	ValidateAccessToken(username string, tokenString string) (bool, error)
	ValidateRefreshToken(tokenString string) (bool, string, error)
	IsIgnoredPath(path string) bool
}

var instance *jwtWrapper
var once sync.Once

type jwtWrapper struct {
	secretKey   string
	ignorePaths []string
}

// Singleton
func NewJWTWrapper(secret string, ignore []string) *jwtWrapper {
	once.Do(func() {
		instance = &jwtWrapper{secretKey: secret, ignorePaths: ignore}
	})
	return instance
}

func GetJWTWrapper() *jwtWrapper {
	return instance
}

func (wrap *jwtWrapper) GenerateJWT(username string, expire int64) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Second * time.Duration(expire)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(wrap.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (wrap *jwtWrapper) ValidateAccessToken(username string, tokenString string) (bool, error) {
	if len(tokenString) <= 0 {
		return false, nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(wrap.secretKey), nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid && token.Claims.(jwt.MapClaims)["username"] == username, nil
}

func (wrap *jwtWrapper) ValidateRefreshToken(tokenString string) (bool, string, error) {
	if len(tokenString) <= 0 {
		return false, "", nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(wrap.secretKey), nil
	})

	if err != nil {
		return false, "", err
	}

	return token.Valid, token.Claims.(jwt.MapClaims)["username"].(string), nil
}

func (wrap *jwtWrapper) IsIgnoredPath(path string) bool {
	return util.InArray(wrap.ignorePaths, path)
}
