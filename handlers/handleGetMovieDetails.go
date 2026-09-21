package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) handleGetMovieDetails(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("movieId")
	movieId, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie id", err)
		return
	}

	movie, err := cfg.DB.GetMovieDetails(r.Context(), movieId)
	if err != err {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, "Movie not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error fetching movie details", err)
		return
	}

	// using this work around to change the show_times to json.RawMessage because sqlc failed to override it from []byte
	response := struct {
		database.GetMovieDetailsRow
		ShowTimes json.RawMessage `json:"show_times"`
	}{
		GetMovieDetailsRow: movie,
		ShowTimes:          movie.ShowTimes,
	}

	respondWithJSON(w, http.StatusOK, response)
}
