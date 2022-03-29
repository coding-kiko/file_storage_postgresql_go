package auth

import (
	//std lib
	"errors"
	"strings"

	// Third party
	"golang.org/x/crypto/bcrypt"
)

func (rd *redisDB) Register(creds Credentials) error {
	if creds.Pwd == "" || !strings.Contains(creds.Email, "@") {
		return errors.New("invalid credentials")
	}
	pwdHsh, err := bcrypt.GenerateFromPassword([]byte(creds.Pwd), 14)
	if err != nil {
		return err
	}
	_, err = rd.conn.Do("HSET", "users:"+creds.Email, "pwd", pwdHsh, "requests", "0")
	if err != nil {
		return err
	}
	return nil
}
