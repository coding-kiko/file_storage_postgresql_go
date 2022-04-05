package file_processing_service

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	updateFileQuery    = `UPDATE files SET file=$1 WHERE id=$2`
	insertRecordsQuery = `INSERT INTO records(file_id, created_at, username, processed)
						  Values($1, $2, $3, $4)`
	getFileQuery    = `SELECT file FROM files WHERE id='%s'`
	getUsernameById = `SELECT username FROM records WHERE file_id='%s'`
)

type postgresDB struct {
	db *sql.DB
}

type Repository interface {
	GetFile(fileId string) ([]byte, string, error)
	CreateFile(form UpdateFileForm) error
}

func NewRepo(db *sql.DB) Repository {
	return &postgresDB{db: db}
}

func (pgDB *postgresDB) GetFile(fileId string) ([]byte, string, error) {
	var data []byte
	var username string
	getFileQuery := fmt.Sprintf(getFileQuery, fileId)
	getUsernameQuery := fmt.Sprintf(getUsernameById, fileId)

	err := pgDB.db.QueryRow(getFileQuery).Scan(&data)
	if err != nil {
		return data, "", errors.New("file not found")
	}
	err = pgDB.db.QueryRow(getUsernameQuery).Scan(&username)
	if err != nil {
		return data, "", errors.New("user not found")
	}
	return data, username, nil
}

func (pgDB *postgresDB) CreateFile(form UpdateFileForm) error {
	// update resized file
	_, err := pgDB.db.Exec(updateFileQuery, form.data, form.id)
	if err != nil {
		return err
	}
	// create new record with file process
	_, err = pgDB.db.Exec(insertRecordsQuery, form.id, form.created_at, form.username, form.processed)
	if err != nil {
		return err
	}
	return nil
}
