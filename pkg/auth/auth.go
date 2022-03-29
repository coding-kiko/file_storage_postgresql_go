package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
)

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
