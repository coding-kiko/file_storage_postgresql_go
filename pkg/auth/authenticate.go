package auth

import (
	//std lib
	"encoding/json"
	"errors"
	"io/ioutil"

	// Third party
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(creds Credentials) (string, error) {
	var users Users
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return "", err
	}

	json.Unmarshal(data, &users)
	for _, u := range users.Users {
		if u.Email == creds.Email {
			if err = bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(creds.Pwd)); err == nil {
				return NewToken(creds.Email), nil
			}
		}
	}
	return "", errors.New("invalid credentials")
}
