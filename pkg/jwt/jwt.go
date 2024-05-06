package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTWrapper interface {
	GenerateJWT(username string) (string, error)
	ValidateToken(username string, tokenString string) (bool, error)
}

type jwtWrapper struct {
	secretKey string
}

func NewJWTWrapper(secret string) *jwtWrapper {
	return &jwtWrapper{secretKey: secret}
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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(wrap.secretKey), nil
	})

	//Token not valid error
	if err != nil {
		return false, err
	}

	return token.Valid && token.Claims.(jwt.MapClaims)["username"] == username, nil
}
