package service

import (
	// std lib
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	// Internal
	"github.com/coding-kiko/file_storage_testing/pkg/files_service/repository"
)

const (
	MAX_UPLOAD_SIZE = 1024 * 1024 * 2 // (2 MB)
)

type service struct {
	repo repository.Respository
}

type ImageService interface {
	GetFile(w http.ResponseWriter, r *http.Request, filename string) error
	CreateFile(w http.ResponseWriter, r *http.Request) error
}

func NewService(repo repository.Respository) ImageService {
	return &service{repo: repo}
}

func (s *service) GetFile(w http.ResponseWriter, r *http.Request, filename string) error {
	filepath, err := s.repo.GetFile(filename)
	defer os.Remove(filepath)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filepath)
	return nil
}

func (s *service) CreateFile(w http.ResponseWriter, r *http.Request) error {
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
	if ext := handler.Filename[len(handler.Filename)-4:]; ext != ".jpg" && ext != ".png" {
		return errors.New("file must be an image")
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

	err = s.repo.CreateFile(data, handler.Filename)
	if err != nil {
		return err
	}
	return nil
}
