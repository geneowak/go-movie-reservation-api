package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleGetMovies(w http.ResponseWriter, r *http.Request) {
	genreParam := r.URL.Query().Get("genre")
	timeParam := r.URL.Query().Get("time")
	var timeQuery time.Time
	if timeParam != "" {
		result, err := time.Parse("15:00", timeParam)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid time param", err)
			return
		}
		timeQuery = result
	}
	dateParam := r.URL.Query().Get("date")
	var dateQuery time.Time
	if dateParam != "" {
		result, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid date query param", err)
			return
		}
		dateQuery = result
	}
	movies, err := cfg.DB.GetShowingMovies(r.Context(), database.GetShowingMoviesParams{
		Genre: &genreParam,
		Time:  &timeQuery,
		Date:  &dateQuery,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "No showing movies at the moment", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching movies", err)
		return
	}

	respondWithJSON(w, http.StatusOK, movies)
}
