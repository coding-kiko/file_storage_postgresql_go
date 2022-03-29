package repository

import (
	//std lib
	"database/sql"
	"net/http"
)

// relative to where the main program is being executed
const (
	MAX_UPLOAD_SIZE = 1024 * 1024 * 2 // (2 MB)
)

type Respository interface {
	GetFile(w http.ResponseWriter, r *http.Request, filename string) error
	CreateFile(w http.ResponseWriter, r *http.Request) error
}

func NewRepo(db *sql.DB) Respository {
	return &postgresDB{db: db}
}
