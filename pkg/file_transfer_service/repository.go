package file_transfer_service

import (
	//std lib
	"database/sql"
	"errors"
	"fmt"
)

var (
	insertFileQuery = `INSERT INTO files(id, name, file)
				  	   Values($1, $2, $3)`
	insertRecordsQuery = `INSERT INTO records(file_id, created_at, username, processed)
						  Values($1, $2, $3, $4)`
	getFileQuery = `SELECT file FROM files WHERE name='%s'`
)

type postgresDB struct {
	db *sql.DB
}

type Repository interface {
	GetFile(filename string) ([]byte, error)
	CreateFile(form CreateFileForm) error
}

func NewRepo(db *sql.DB) Repository {
	return &postgresDB{db: db}
}

func (pgDB *postgresDB) GetFile(filename string) ([]byte, error) {
	var data []byte
	getFileQuery := fmt.Sprintf(getFileQuery, filename)

	err := pgDB.db.QueryRow(getFileQuery).Scan(&data)
	if err != nil {
		return data, errors.New("file not found")
	}
	return data, nil
}

func (pgDB *postgresDB) CreateFile(form CreateFileForm) error {
	_, err := pgDB.db.Exec(insertFileQuery, form.id, form.filename, form.data)
	if err != nil {
		return err
	}
	_, err = pgDB.db.Exec(insertRecordsQuery, form.id, form.created_at, form.username, form.processed)
	if err != nil {
		return err
	}
	return nil
}
