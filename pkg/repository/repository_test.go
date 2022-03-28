package repository

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetFile(w http.ResponseWriter, r *http.Request, filename string) error {
	args := m.Called(w, r, filename)
	return args.Error(0)
}

func (m *MockRepository) CreateFile(w http.ResponseWriter, r *http.Request) error {
	args := m.Called(w, r)
	return args.Error(0)
}

func TestDownloadFileNotFound(t *testing.T) {
	var w http.ResponseWriter
	var r *http.Request

	fileNotFoundErrMsg := "file not found"
	mockedRepo := new(MockRepository)
	mockedRepo.On("GetFile", mock.Anything, mock.Anything, "notFound.pdf").Return(errors.New(fileNotFoundErrMsg))
	err := mockedRepo.GetFile(w, r, "notFound.pdf")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fileNotFoundErrMsg)
}

func TestFileSizeExceeded(t *testing.T) {
	var w http.ResponseWriter
	var r *http.Request
	sizeExceededMsg := "file size exceeded"

	mockedRepo := new(MockRepository)
	mockedRepo.On("CreateFile", mock.Anything, r).Return(errors.New(sizeExceededMsg))
	err := mockedRepo.CreateFile(w, r)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), sizeExceededMsg)
}

func TestDuplicateFile(t *testing.T) {
	var w http.ResponseWriter
	var r *http.Request
	duplicateFileMsg := "error uploading file"

	mockedRepo := new(MockRepository)
	mockedRepo.On("CreateFile", mock.Anything, mock.Anything).Return(errors.New(duplicateFileMsg))
	err := mockedRepo.CreateFile(w, r)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), duplicateFileMsg)
}
