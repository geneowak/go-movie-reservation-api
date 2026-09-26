package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/geneowak/cinehold/internal/database"
	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleGetMovies(w http.ResponseWriter, r *http.Request) {
	// The optional filters are sqlc.narg params, so an absent filter has to reach
	// the query as NULL rather than as a pointer to a zero value. A pointer to
	// "0001-01-01" would make the date guard `start_date <= $date` false for every
	// showtime and the whole listing would come back empty.
	params := database.GetShowingMoviesParams{}

	if genreParam := r.URL.Query().Get("genre"); genreParam != "" {
		params.Genre = &genreParam
	}

	if timeParam := r.URL.Query().Get("time"); timeParam != "" {
		timeQuery, err := time.Parse("15:00", timeParam)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid time param", err)
			return
		}
		params.Time = &timeQuery
	}

	if dateParam := r.URL.Query().Get("date"); dateParam != "" {
		dateQuery, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid date query param", err)
			return
		}
		params.Date = &dateQuery
	}

	movies, err := cfg.DB.GetShowingMovies(r.Context(), params)
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
