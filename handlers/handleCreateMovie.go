package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/database"
)

type createMovieRequest struct {
	Name            string   `json:"name" validate:"required,alphanum,min=3"`
	Description     string   `json:"description" validate:"required,min=10"`
	DurationInMins  int      `json:"duration_in_mins" validate:"required,number"`
	TrailerUrl      string   `json:"trailer_url" validate:"required,url"`
	Genre           []string `json:"genre" validate:"required,min=1,dive,required"`
	PgRating        string   `json:"pg_rating" validate:"required,min=2"`
	ExperienceTypes []string `json:"experience_types" validate:"required,min=1,dive,required"`
}

func (cfg *ApiConfig) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	var req createMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleJsonDecodeError(w, err)
		return
	}

	if err := cfg.Validate.Struct(req); err != nil {
		handleValidationErrors(w, err)
		return
	}

	movie, err := cfg.DB.CreateMovie(r.Context(), database.CreateMovieParams{
		Name:            req.Name,
		Description:     req.Description,
		DurationInMins:  int32(req.DurationInMins),
		TrailerUrl:      req.TrailerUrl,
		Genre:           req.Genre,
		PgRating:        req.PgRating,
		ExperienceTypes: req.ExperienceTypes,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving movie details.", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, movie)
}
