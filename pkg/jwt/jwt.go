package jwt

import (
	"sync"
	"time"

	"github.com/Beretta350/authentication/pkg/util"
	"github.com/golang-jwt/jwt"
)

type JWTWrapper interface {
	GenerateJWT(username string) (string, error)
	ValidateToken(username string, tokenString string) (bool, error)
	IsIgnoredPath(path string) bool
	SetSecretKey(secret string)
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

func (wrap *jwtWrapper) GenerateJWT(username string) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(wrap.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (wrap *jwtWrapper) ValidateToken(username string, tokenString string) (bool, error) {
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

func (wrap *jwtWrapper) IsIgnoredPath(path string) bool {
	return util.InArray(wrap.ignorePaths, path)
}

func (wrap *jwtWrapper) SetSecretKey(secret string) {
	wrap.secretKey = secret
}
