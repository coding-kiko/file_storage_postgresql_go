package repository

import (
	// std lib
	"errors"
	"io"
	"net/http"
	"os"
)

type fileStorage struct {
	relativePath string
}

// Get file from local storage directory "~/Desktop/repos/file_storage_testing/repository/file_storage"
func (fs *fileStorage) GetFile(w http.ResponseWriter, r *http.Request, filename string) error {
	filepath := fs.relativePath + filename // NOTICE: path is relative to the program execution path
	// check if file exists
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return errors.New("file not found")
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filepath) // NOTE: This function already sets header status code to 200
	return nil
}

// Create new file in local storage directory "~/Desktop/repos/file_storage_testing/repository/file_storage"
func (fs *fileStorage) CreateFile(w http.ResponseWriter, r *http.Request) error {
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
	defer file.Close()
	f, err := os.OpenFile(fs.relativePath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("error uploading file")
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}
