package auth

import (
	// std lib
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	// Third party
	"github.com/dgrijalva/jwt-go"
)

var (
	secretKey = "secret"
	issuer    = "localhost:5000/"
)

func ValidateJwt(token string) error {
	tk, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return errors.New("invalid token")
	}
	claims, ok := tk.Claims.(*Claims)
	if !ok {
		return errors.New("invalid token")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return errors.New("expired token")
	}
	err = CheckAvailableRequests(claims.Name)
	return err
}

// open cache file and check if user has available requests with current jwt
func CheckAvailableRequests(email string) error {
	var users Users
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return err
	}
	json.Unmarshal(data, &users)
	for i, u := range users.Users {
		if u.Email == email {
			if users.Users[i].Requests == 0 {
				return errors.New("request amount exceeded, sign in again")
			}
			users.Users[i].Requests = u.Requests - 1
			data, err = json.MarshalIndent(users, "", "\t")
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(cacheFile, data, 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

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
