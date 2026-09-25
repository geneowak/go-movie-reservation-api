package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type updateMovieRequest struct {
	Name            string   `json:"name" validate:"required,min=3"`
	Description     string   `json:"description" validate:"required,min=10"`
	DurationInMins  int      `json:"duration_in_mins" validate:"required,number"`
	TrailerUrl      string   `json:"trailer_url" validate:"required,url"`
	Genre           []string `json:"genre" validate:"required,min=1,dive,required"`
	PgRating        string   `json:"pg_rating" validate:"required,min=2"`
	ExperienceTypes []string `json:"experience_types" validate:"required,min=1,dive,required"`
}

func (cfg *ApiConfig) handleUpdateMovieDetails(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("movieId")
	movieId, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid movie id", err)
		return
	}

	var req updateMovieRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleJsonDecodeError(w, err)
		return
	}

	if err = cfg.Validate.Struct(req); err != nil {
		handleValidationErrors(w, err)
		return
	}

	updatedMovie, err := cfg.DB.UpdateMovieDetails(r.Context(), database.UpdateMovieDetailsParams{
		Name:            req.Name,
		Description:     req.Description,
		DurationInMins:  int32(req.DurationInMins),
		TrailerUrl:      req.TrailerUrl,
		Genre:           req.Genre,
		PgRating:        req.PgRating,
		ExperienceTypes: req.ExperienceTypes,
		ID:              movieId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Movie not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error saving movie details", err)
		return
	}

	respondWithJSON(w, http.StatusOK, updatedMovie)
}
