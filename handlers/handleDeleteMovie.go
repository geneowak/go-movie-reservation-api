package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/geneowak/cinehold/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("movieId")
	movieId, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie id", err)
		return
	}
	// a movie that has an ongoing booking must not be deletable because it could have bookings
	movieCurrentShowTimes, err := cfg.DB.GetMovieOngoingShowtimes(r.Context(), database.GetMovieOngoingShowtimesParams{
		ID:      movieId,
		EndDate: time.Now(),
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Movie not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error saving checking movie showtimes", err)
		return
	}
	var currentShowTimes []json.RawMessage
	_ = json.Unmarshal(movieCurrentShowTimes.ShowTimes, &currentShowTimes)
	if len(currentShowTimes) > 0 {
		respondWithError(w, http.StatusBadRequest, "Cannot delete a movie with running showtimes", errors.New("Movie has running show times"))
		return
	}

	err = cfg.DB.DeleteMovie(r.Context(), movieId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting movie", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
