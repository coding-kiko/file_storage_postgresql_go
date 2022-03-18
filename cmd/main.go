package main

import (
	// std lib
	"fmt"
	"log"
	"net/http"

	// third party
	"github.com/caarlos0/env"

	// Internal
	"github.com/coding-kiko/file_storage_testing/repository"
	"github.com/coding-kiko/file_storage_testing/server"
)

func main() {
	// Parse env variables to cfg struct
	cfg := Config{}
	env.Parse(&cfg)

	dbCfg := repository.DBConfig{DB: cfg.DB}
	db := repository.NewRepo(dbCfg)

	mux := http.NewServeMux()
	mux.Handle("/file", server.GetFileHandler(db))
	mux.Handle("/upload", server.CreateFileHandler(db))
	fmt.Printf("Started listening on %d\n", cfg.Port)
	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	err := http.ListenAndServe(addr, mux)
	log.Fatal(err)
}

type Config struct {
	Port int    `env:"PORT" envDefault:"3000"`
	DB   string `env:"DB" envDefault:"fs"`
}
