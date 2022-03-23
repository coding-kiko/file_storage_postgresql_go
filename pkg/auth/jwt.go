package auth

import (
	// std lib
	"time"

	// Third party
	"github.com/dgrijalva/jwt-go"
)

var (
	secretKey = "secret"
	issuer    = "localhost:5000/"
)

func NewToken(name string) string {
	signingKey := []byte(secretKey)
	claims := Claims{
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(5) * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(signingKey)
	return signedToken
}

type Claims struct {
	Name string
	jwt.StandardClaims
}

func (c Claims) Valid() error {
	return nil
}
