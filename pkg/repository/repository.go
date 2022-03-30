package repository

import (
	//std lib
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
)

var (
	insertFileQuery = `INSERT INTO files(name, file)
				  	   Values($1, $2)`
	getFileQuery = `SELECT file FROM files WHERE name='%s'`
)

type postgresDB struct {
	db *sql.DB
}

type Respository interface {
	GetFile(filename string) (string, error)
	CreateFile(binaryData []byte, filename string) error
}

func NewRepo(db *sql.DB) Respository {
	return &postgresDB{db: db}
}

func (pgDB *postgresDB) GetFile(filename string) (string, error) {
	var data []byte
	filepath := "./tmp/" + filename
	getFileQuery := fmt.Sprintf(getFileQuery, filename)

	err := pgDB.db.QueryRow(getFileQuery).Scan(&data)
	if err != nil {
		return "", errors.New("file not found")
	}
	err = ioutil.WriteFile(filepath, data, 0666)
	if err != nil {
		return "", errors.New("error writing file")
	}
	return filepath, nil
}

func (pgDB *postgresDB) CreateFile(binaryData []byte, filename string) error {
	_, err := pgDB.db.ExecContext(context.Background(), insertFileQuery, filename, binaryData)
	if err != nil {
		// probably duplicate (filename is primary key)
		return errors.New("file already exists")
	}
	return nil
}
