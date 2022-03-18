package repository

import (
	"net/http"
)

// relative to where the main program is being executed
const (
	fileStorageRelativePath = "./repository/file_storage/"
	MAX_UPLOAD_SIZE         = 1024 * 1024 * 2 // (2 MB)
)

type Respository interface {
	GetFile(w http.ResponseWriter, r *http.Request, filename string) error
	CreateFile(w http.ResponseWriter, r *http.Request) error
}

func NewRepo(cfg DBConfig) Respository {
	// check if file_storage is chosen as database (default)
	if cfg.DB == "fs" {
		return &fileStorage{relativePath: fileStorageRelativePath}
	}
	return &postgresDB{}
}

type DBConfig struct {
	DB string
}
