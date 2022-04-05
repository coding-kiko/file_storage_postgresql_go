package file_transfer_service

import (
	// std lib
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	// Third party
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const (
	MAX_UPLOAD_SIZE = 1024 * 1024 * 3 // (3 MB)
)

type service struct {
	repo       Repository
	rabbitConn amqp.Connection
}

type ImageService interface {
	GetFile(w http.ResponseWriter, r *http.Request, filename string) error
	CreateFile(w http.ResponseWriter, r *http.Request, form string) error
}

func NewService(repo Repository, rbConn amqp.Connection) ImageService {
	return &service{repo: repo, rabbitConn: rbConn}
}

func (s *service) GetFile(w http.ResponseWriter, r *http.Request, filename string) error {
	byteData, err := s.repo.GetFile(filename)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(byteData))
	return nil
}

func (s *service) CreateFile(w http.ResponseWriter, r *http.Request, username string) error {
	var byteData []byte

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

	// parse file to byte slice
	byteData, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fileId := uuid.New().String()
	form := CreateFileForm{
		id:         fileId,
		data:       byteData,
		username:   username,
		filename:   handler.Filename,
		processed:  false,
		created_at: time.Now().Format("2006-01-02 15:04"),
	}
	err = s.repo.CreateFile(form)
	if err != nil {
		return err
	}
	err = sendFileIdToRabbitQueue(fileId, s.rabbitConn)
	if err != nil {
		return err
	}
	return nil
}

func sendFileIdToRabbitQueue(fileId string, conn amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return errors.New("failed to open channel")
	}
	defer ch.Close()

	// We create a Queue to send the message to.
	q, err := ch.QueueDeclare(
		"new-file-id-queue", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return errors.New("failed to declare a queue")
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fileId),
		})
	if err != nil {
		return errors.New("failed to publish message")
	}
	return nil
}

type CreateFileForm struct {
	data       []byte
	filename   string
	username   string
	created_at string
	processed  bool
	id         string
}
