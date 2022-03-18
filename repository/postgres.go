package repository

import "net/http"

type postgresDB struct {
}

func (db *postgresDB) GetFile(w http.ResponseWriter, r *http.Request, filename string) error {

	return nil
}

func (db *postgresDB) CreateFile(w http.ResponseWriter, r *http.Request) error {
	return nil
}
