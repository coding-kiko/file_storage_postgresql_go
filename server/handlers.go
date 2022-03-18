package server

import (
	"net/http"

	"github.com/coding-kiko/file_storage_testing/repository"
)

func GetFileHandler(db repository.Respository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// check if filename is not passed correctly through query params
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status": 400, "message": "Bad Request"}`))
			return
		}
		// check if another method other than GET has reached the endpoint
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"status": 405, "message": "Method not allowed"}`))
			return
		}
		err := db.GetFile(w, r, filename)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"status": 404, "message": "File not found"}`))
			return
		}
	})
}

func CreateFileHandler(db repository.Respository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// check if another method other than GET has reached the endpoint
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"status": 405, "message": "Method not allowed"}`))
			return
		}
		err := db.CreateFile(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status": 400, "message": "File size is too large"}`))
			return
		}
	})
}
