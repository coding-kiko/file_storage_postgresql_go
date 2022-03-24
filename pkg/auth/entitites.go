package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Email    string `json:"email"`
	Pwd      string `json:"pwd"`
	Requests int    `json:"requests"`
}

type Credentials struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

type Claims struct {
	Name string
	jwt.StandardClaims
}

func (c Claims) Valid() error {
	return nil
}
