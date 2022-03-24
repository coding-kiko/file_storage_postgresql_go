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
	for i, u := range users.Users {
		if u.Email == creds.Email {
			if err = bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(creds.Pwd)); err == nil {
				// Sets 5 requests for the user to use before expiring
				users.Users[i].Requests = 5
				data, err = json.MarshalIndent(users, "", "\t")
				if err != nil {
					return "", err
				}
				err = ioutil.WriteFile(cacheFile, data, 0644)
				if err != nil {
					return "", err
				}
				return NewToken(creds.Email), nil
			}
		}
	}
	return "", errors.New("invalid credentials")
}
