package auth

import (
	//std lib
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	// Third party
	"golang.org/x/crypto/bcrypt"
)

const (
	cacheFile = "./cache/users.json"
)

func Register(creds Credentials) error {
	var users Users
	if creds.Pwd == "" || !strings.Contains(creds.Email, "@") {
		return errors.New("invalid credentials")
	}

	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return err
	}

	json.Unmarshal(data, &users)
	pwdHsh, err := bcrypt.GenerateFromPassword([]byte(creds.Pwd), 14)
	if err != nil {
		return err
	}
	newUser := User{Email: creds.Email, Pwd: string(pwdHsh)}
	users.Users = append(users.Users, newUser)

	newUsersList := Users{Users: users.Users}
	data, err = json.MarshalIndent(newUsersList, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cacheFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
