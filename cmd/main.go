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
	"github.com/coding-kiko/file_storage_testing/pkg/file_processing_service"
	"github.com/coding-kiko/file_storage_testing/pkg/file_transfer_service"
	"github.com/coding-kiko/file_storage_testing/pkg/server"

	// third party
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
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

	// redis connection
	redisConn, err := redis.Dial("tcp", *redisAddr)
	if err != nil {
		panic(err)
	}
	defer redisConn.Close()

	// make postgres connection
	connString := fmt.Sprintf("postgres://postgres:%s@%s/%s%s", *pwd, *dBAddr, *database, "?sslmode=disable")
	postgresDb, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	defer postgresDb.Close()

	// rabbitmq connection
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer rabbitConn.Close()

	// initialize file processing service layers
	fileProcessingRepo := file_processing_service.NewRepo(postgresDb)
	fileProcessingServ := file_processing_service.NewService(fileProcessingRepo)
	go file_processing_service.NewQueueConsumer(*rabbitConn, fileProcessingServ)

	// initialize file transfer service layers
	fileTransferRepo := file_transfer_service.NewRepo(postgresDb)
	fileTransferServ := file_transfer_service.NewService(fileTransferRepo, *rabbitConn)
	serviceHandlers := server.NewServiceHandlers(fileTransferServ)

	// initialize authentication service layers
	authRepo := auth.NewRedisRepo(redisConn)
	authHandlers := server.NewAuthHandlers(authRepo)

	mux := server.Init(serviceHandlers, authHandlers)
	fmt.Printf("Started listening on %s\n", *port)
	addr := fmt.Sprintf("localhost:%s", *port)
	err = http.ListenAndServe(addr, mux)
	log.Fatal(err)
}
