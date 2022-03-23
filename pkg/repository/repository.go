package repository

import (
	//std lib
	"database/sql"
	"net/http"
)

// relative to where the main program is being executed
const (
	fileStorageRelativePath = "./pkg/repository/file_storage/"
	MAX_UPLOAD_SIZE         = 1024 * 1024 * 2 // (2 MB)
)

type Respository interface {
	GetFile(w http.ResponseWriter, r *http.Request, filename string) error
	CreateFile(w http.ResponseWriter, r *http.Request) error
}

func NewRepo(cfg DBConfig) Respository {
	// check if file_storage is chosen as database (default)
	if cfg.Storage == "fs" {
		return &fileStorage{relativePath: fileStorageRelativePath}
	}
	return &postgresDB{db: cfg.DB}
}

type DBConfig struct {
	Storage string
	DB      *sql.DB
}
