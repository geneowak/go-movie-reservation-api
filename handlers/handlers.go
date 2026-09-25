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

	// Admin routes
	mux.HandleFunc("POST /api/movies", cfg.middlewareAdminAuth(cfg.handleCreateMovie))
	mux.HandleFunc("POST /api/movies/{movieId}/show-times", cfg.middlewareAdminAuth(cfg.handleCreateMovieShowTime))
	mux.HandleFunc("GET /api/admin/movies", cfg.middlewareAdminAuth(cfg.handleAdminShowMovies))

	mux.HandleFunc("POST /api/locations", cfg.middlewareAdminAuth(cfg.handleCreateLocation))

	mux.HandleFunc("POST /api/cinemas", cfg.middlewareAdminAuth(cfg.handleCreateCinema))

	// user routes
	mux.HandleFunc("GET /api/movies", cfg.middlewareAuth(cfg.handleGetMovies))
	mux.HandleFunc("GET /api/movies/{movieId}", cfg.middlewareAuth(cfg.handleGetMovieDetails))

	mux.HandleFunc("POST /api/seats/reserve", cfg.middlewareAuth(cfg.handleReserveSeat))
	mux.HandleFunc("POST /api/seats/book", cfg.middlewareAuth(cfg.handleSeatBooking))
	mux.HandleFunc("GET /api/bookings", cfg.middlewareAuth(cfg.handleGetBookings))

	return &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
