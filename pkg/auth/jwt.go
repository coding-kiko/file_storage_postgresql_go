package auth

import (
	// std lib
	"errors"
	"fmt"
	"time"

	// Third party
	"github.com/dgrijalva/jwt-go"
)

var (
	secretKey = "secret"
	issuer    = "localhost:5000/"
)

type Claims struct {
	Name string
	jwt.StandardClaims
}

func (c Claims) Valid() error {
	return nil
}

func (rd *redisDB) ValidateJwt(token string) (string, error) {
	// parse token with secret passphrase
	tk, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", errors.New("invalid token")
	}
	claims, ok := tk.Claims.(*Claims)
	if !ok {
		return "", errors.New("invalid token")
	}
	// check expiry date
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return "", errors.New("expired token")
	}

	// check and rate limit the number of requests
	req, err := rd.conn.Do("HGET", "users:"+claims.Name, "requests")
	if err != nil {
		return "", err
	}

	n := fmt.Sprintf("%s", req)
	if n == "0" {
		return "", errors.New("request amount exceeded, sign in again")
	}
	// subtract "1" to the current number of requests
	n = string(rune(int([]rune(n)[0]) - 1))

	_, err = rd.conn.Do("HSET", "users:"+claims.Name, "requests", n)
	if err != nil {
		return "", err
	}
	return claims.Name, nil
}

func NewToken(name string) string {
	signingKey := []byte(secretKey)
	claims := Claims{
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(15) * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(signingKey)
	return signedToken
}
