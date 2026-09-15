package handlers

import (
	"net/http"
	"time"
)

func SetupServer(cfg *ApiConfig, filePathRoot, port string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/signup", cfg.handleCreateUser)

	mux.HandleFunc("POST /api/login", cfg.handleLogin)

	mux.HandleFunc("POST /api/auth/refresh", cfg.handleRefreshToken)
	mux.HandleFunc("POST /api/auth/revoke", cfg.handleRevokeToken)

	mux.HandleFunc("POST /api/movies", cfg.middlewareAuth(cfg.handleCreateMovie))

	mux.HandleFunc("POST /api/locations", cfg.middlewareAuth(cfg.handleCreateLocation))

	mux.HandleFunc("POST /api/cinemas", cfg.middlewareAuth(cfg.handleCreateCinema))

	return &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
