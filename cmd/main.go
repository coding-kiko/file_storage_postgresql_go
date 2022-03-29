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

	rconn, err := redis.Dial("tcp", *redisAddr)
	if err != nil {
		panic(err)
	}
	defer rconn.Close()
	rd := auth.NewRedisRepo(rconn)

	// make postgres connection
	connString := fmt.Sprintf("postgres://postgres:%s@%s/%s%s", *pwd, *dBAddr, *database, "?sslmode=disable")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	pg := repository.NewRepo(db)

	mux := http.NewServeMux()
	mux.Handle("/file", server.JwtMiddleware(server.GetFileHandler(pg), rd))
	mux.Handle("/upload", server.JwtMiddleware(server.CreateFileHandler(pg), rd))
	mux.Handle("/authenticate", server.AuthenticateHandler(rd))
	mux.Handle("/register", server.RegisterHandler(rd))

	fmt.Printf("Started listening on %s\n", *port)
	addr := fmt.Sprintf("localhost:%s", *port)
	err = http.ListenAndServe(addr, mux)
	log.Fatal(err)
}
