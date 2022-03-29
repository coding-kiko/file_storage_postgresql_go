package auth

import (
	//std lib

	"errors"
	"fmt"

	// Third party
	"golang.org/x/crypto/bcrypt"
)

func (rd *redisDB) Authenticate(creds Credentials) (string, error) {
	resp, err := rd.conn.Do("HGET", "users:"+creds.Email, "pwd")
	if err != nil || resp == nil {
		return "", errors.New("invalid credentials")
	}
	pwdHsh := fmt.Sprintf("%s", resp)
	if err = bcrypt.CompareHashAndPassword([]byte(pwdHsh), []byte(creds.Pwd)); err != nil {
		return "", errors.New("invalid credentials")
	}
	_, err = rd.conn.Do("HSET", "users:"+creds.Email, "requests", "5")
	if err != nil {
		return "", err
	}
	return NewToken(creds.Email), nil
}
