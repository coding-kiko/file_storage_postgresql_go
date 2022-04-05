package server

import (
	// std lib
	"net/http"
)

func Init(sh ServiceHandlers, ah AuthHandlers) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/download", ah.JwtMiddleware(sh.GetFileHandler()))
	mux.Handle("/upload", ah.JwtMiddleware(sh.CreateFileHandler()))
	mux.Handle("/authenticate", ah.AuthenticateHandler())
	mux.Handle("/register", ah.RegisterHandler())
	return mux
}
