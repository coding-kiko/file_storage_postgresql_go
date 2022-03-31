package auth

import (
	// Std lib
	"errors"
	"fmt"
	"strings"

	// Third party
	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

type redisDB struct {
	conn redis.Conn
}

type RedisRespository interface {
	Authenticate(creds Credentials) (string, error)
	Register(creds Credentials) error
	ValidateJwt(token string) error
}

func NewRedisRepo(conn redis.Conn) RedisRespository {
	return &redisDB{conn: conn}
}

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
