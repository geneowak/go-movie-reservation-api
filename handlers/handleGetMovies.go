package handlers

import (
	"net/http"
	"strings"
)

func (cfg *ApiConfig) handleGetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := cfg.DB.GetShowingMovies(r.Context())
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, "No showing movies at the moment", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching movies", err)
		return
	}

	respondWithJSON(w, http.StatusOK, movies)
}
