package main

import (
	// std lib
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	// Internal
	"github.com/coding-kiko/file_storage_testing/pkg/auth"
	"github.com/coding-kiko/file_storage_testing/pkg/repository"
	"github.com/coding-kiko/file_storage_testing/pkg/server"
	"github.com/coding-kiko/file_storage_testing/pkg/service"

	// third party
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
)

var (
	port      = flag.String("PORT", "5000", "api port")
	dBAddr    = flag.String("DBADDR", "localhost:5432", "database address")
	pwd       = flag.String("PWD", "admin", "postgres password")
	database  = flag.String("DB", "file_storage", "postgres database")
	redisAddr = flag.String("REDIS", "localhost:6379", "redis db addr")
)

func main() {
	flag.Parse()

	redisConn, err := redis.Dial("tcp", *redisAddr)
	if err != nil {
		panic(err)
	}
	defer redisConn.Close()

	// make postgres connection
	connString := fmt.Sprintf("postgres://postgres:%s@%s/%s%s", *pwd, *dBAddr, *database, "?sslmode=disable")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	serviceHandlers := server.NewServiceHandlers(service.NewService(repository.NewRepo(db)))
	authHandlers := server.NewAuthHandlers(auth.NewRedisRepo(redisConn))

	mux := server.Init(serviceHandlers, authHandlers)
	fmt.Printf("Started listening on %s\n", *port)
	addr := fmt.Sprintf("localhost:%s", *port)
	err = http.ListenAndServe(addr, mux)
	log.Fatal(err)
}
