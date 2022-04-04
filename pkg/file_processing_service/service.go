package file_processing_service

import (
	// std lib
	"bytes"
	"image/jpeg"
	"time"

	// Third party
	"github.com/nfnt/resize"
)

type service struct {
	repo Repository
}

type FileProcessingService interface {
	ResizeImage(fileId string) error
}

func NewService(repo Repository) FileProcessingService {
	return &service{repo: repo}
}

func (s *service) ResizeImage(fileId string) error {
	byteData, username, err := s.repo.GetFile(fileId)
	if err != nil {
		return err
	}

	// resize image
	reader := bytes.NewReader(byteData)
	img, err := jpeg.Decode(reader)
	if err != nil {
		return err
	}
	// Resize to width 1000 and preserve ratio
	m := resize.Resize(1000, 0, img, resize.Lanczos3)
	newBuf := new(bytes.Buffer)
	err = jpeg.Encode(newBuf, m, nil)
	if err != nil {
		return err
	}
	resizedData := newBuf.Bytes()

	form := UpdateFileForm{
		data:       resizedData,
		id:         fileId,
		username:   username,
		processed:  true,
		created_at: time.Now().Format("2006-01-02 15:04"),
	}
	err = s.repo.CreateFile(form)
	if err != nil {
		return err
	}
	return nil
}

type UpdateFileForm struct {
	data       []byte
	username   string
	created_at string
	processed  bool
	id         string
}
