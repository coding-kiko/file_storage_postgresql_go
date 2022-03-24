package repository

import (
	//std lib
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	insertFileQuery = `INSERT INTO files(name, file)
				  	   Values($1, $2)`
	getFileQuery = `SELECT file FROM files WHERE name='%s'`
)

type postgresDB struct {
	db *sql.DB
}

func (pgDB *postgresDB) GetFile(w http.ResponseWriter, r *http.Request, filename string) error {
	var data []byte
	filepath := "./tmp/" + filename
	getFileQuery := fmt.Sprintf(getFileQuery, filename)

	err := pgDB.db.QueryRow(getFileQuery).Scan(&data)
	if err != nil {
		panic(err.Error())
	}
	err = ioutil.WriteFile(filepath, data, 0666)
	if err != nil {
		return errors.New("error writing file")
	}
	defer os.Remove(filepath)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filepath)
	return nil
}

func (pgDB *postgresDB) CreateFile(w http.ResponseWriter, r *http.Request) error {
	// check if file size is within permited range
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		return errors.New("file size exceeded")
	}
	file, handler, err := r.FormFile("filename")
	if err != nil {
		return errors.New("error forming file")
	}
	// TODO: find a way to do this without having to download file to the disk
	f, err := os.OpenFile("./tmp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("error uploading file")
	}
	defer f.Close()
	defer os.Remove("./tmp/" + handler.Filename)
	io.Copy(f, file)
	data, _ := ioutil.ReadFile("./tmp/" + handler.Filename)
	// TODO END

	_, err = pgDB.db.ExecContext(context.Background(), insertFileQuery, handler.Filename, data)
	if err != nil {
		panic(err)
	}
	return nil
}
