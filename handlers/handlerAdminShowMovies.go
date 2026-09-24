package handlers

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleAdminShowMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := cfg.DB.GetAllMovies(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "No movies found.", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch movie details.", err)
		return
	}

	respondWithJSON(w, http.StatusOK, movies)
}
