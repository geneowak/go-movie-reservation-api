package handlers

import (
	"net/http"
	"strings"
)

func (cfg *ApiConfig) handleAdminShowMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := cfg.DB.GetAllMovies(r.Context())
	if err != nil {
		if strings.Contains(err.Error(), EmptyResultSet) {
			respondWithError(w, http.StatusNotFound, "No movies found.", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch movie details.", err)
		return
	}

	respondWithJSON(w, http.StatusOK, movies)
}
