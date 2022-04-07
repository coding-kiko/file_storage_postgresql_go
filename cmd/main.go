package main

import (
	// std lib
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// Internal
	"github.com/coding-kiko/file_storage_testing/pkg/auth"
	"github.com/coding-kiko/file_storage_testing/pkg/file_processing_service"
	"github.com/coding-kiko/file_storage_testing/pkg/file_transfer_service"
	"github.com/coding-kiko/file_storage_testing/pkg/server"

	// third party
	"github.com/caarlos0/env/v6"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type config struct {
	ApiPort           string `env:"API_PORT" envDefault:"5000"`
	PostgresDB        string `env:"POSTGRES_DB" envDefault:"file_storage"`
	PostgresContainer string `env:"POSTGRES:CONTAINER" envDefault:"postgres"`
	PostgresPort      string `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresPwd       string `env:"POSTGRES_PWD" envDefault:"admin"`
	RedisContainer    string `env:"REDIS_CONTAINER" envDefault:"redis"`
	RedisPort         string `env:"REDIS_PORT" envDefault:"6379"`
	RabbitContainer   string `env:"RABBIT_CONTAINER" envDefault:"rabbitmq"`
	RabbitPort        string `env:"RABBIT_PORT" envDefault:"5672"`
}

func main() {
	cfg := config{}
	err := env.Parse(&cfg)

	// redis connection
	redisConnString := fmt.Sprintf("%s:%s", cfg.RedisContainer, cfg.RedisPort)
	redisConn, err := redis.Dial("tcp", redisConnString)
	if err != nil {
		panic(err)
	}
	defer redisConn.Close()

	// make postgres connection
	postgresConnString := fmt.Sprintf("postgres://postgres:%s@%s:%s/%s%s", cfg.PostgresPwd, cfg.PostgresContainer, cfg.PostgresPort, cfg.PostgresDB, "?sslmode=disable")
	postgresDb, err := sql.Open("postgres", postgresConnString)
	if err != nil {
		panic(err)
	}
	defer postgresDb.Close()

	// rabbitmq connection
	rabbitConnString := fmt.Sprintf("amqp://%s:%s/", cfg.RabbitContainer, cfg.RabbitPort)
	rabbitConn, err := amqp.Dial(rabbitConnString)
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
	fmt.Printf("Started listening on %s\n", cfg.ApiPort)
	addr := fmt.Sprintf("0.0.0.0:%s", cfg.ApiPort)
	err = http.ListenAndServe(addr, mux)
	log.Fatal(err)
}
