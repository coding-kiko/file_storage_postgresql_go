package main

import (
	// std lib

	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	// Internal
	"github.com/coding-kiko/file_storage_testing/pkg/repository"
	"github.com/coding-kiko/file_storage_testing/pkg/server"

	// third party
	_ "github.com/lib/pq"
)

var (
	port     = flag.String("PORT", "5000", "api port")
	storage  = flag.String("STORAGE", "fs", "storage type")
	dBAddr   = flag.String("DBADDR", "localhost:5432", "database address")
	pwd      = flag.String("PWD", "admin", "postgres password")
	database = flag.String("DB", "file_storage", "postgres database")
)

func main() {
	flag.Parse()

	dbCfg := repository.DBConfig{Storage: *storage}
	if *storage == "postgres" {
		connString := fmt.Sprintf("postgres://postgres:%s@%s/%s%s", *pwd, *dBAddr, *database, "?sslmode=disable")
		db, err := sql.Open("postgres", connString)
		if err != nil {
			panic(err.Error())
		}
		dbCfg.DB = db
		defer db.Close()
	}
	db := repository.NewRepo(dbCfg)

	mux := http.NewServeMux()
	// mux.Handle("/file", server.JwtMiddleware(server.GetFileHandler(db)))
	// mux.Handle("/upload", server.JwtMiddleware(server.CreateFileHandler(db)))
	mux.Handle("/file", server.GetFileHandler(db))
	mux.Handle("/upload", server.CreateFileHandler(db))
	mux.Handle("/authenticate", server.AuthenticateHandler())
	mux.Handle("/register", server.RegisterHandler())

	fmt.Printf("Started listening on %s\n", *port)
	addr := fmt.Sprintf("localhost:%s", *port)
	err := http.ListenAndServe(addr, mux)
	log.Fatal(err)
}
