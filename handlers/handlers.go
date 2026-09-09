package handlers

import (
	"net/http"
	"time"
)

func SetupServer(cfg *ApiConfig, filePathRoot, port string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/users", cfg.handleCreateUser)

	mux.HandleFunc("POST /api/login", cfg.handleLogin)

	return &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
