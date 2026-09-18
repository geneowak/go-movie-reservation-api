package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
)

type createMovieShowTimeRequest struct {
	StartTime      string  `json:"start_time" validate:"required,datetime=15:04"`
	Price          float32 `json:"price" validate:"required,number,gt=0"`
	Description    string  `json:"description" validate:"required"`
	PriceCurrency  string  `json:"price_currency" validate:"required,iso4217"`
	CinemaID       string  `json:"cinema_id" validate:"required,uuid_rfc4122,cinema-exists"`
	ExperienceType string  `json:"experience_type" validate:"required"`
	StartDate      string  `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate        string  `json:"end_date" validate:"required,datetime=2006-01-02,gtecsfield=StartDate"`
}

func (cfg *ApiConfig) handleCreateMovieShowTime(w http.ResponseWriter, r *http.Request) {
	paramId := r.PathValue("movieId")
	movieId, err := uuid.Parse(paramId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie id.", err)
		return
	}
	// ensure that the movie id that was selected exists.
	if exists, err := cfg.DB.CheckMovieById(r.Context(), movieId); err != nil || !exists {
		respondWithError(w, http.StatusBadRequest, "Invalid movie id.", err)
		return
	}

	var req createMovieShowTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := cfg.Validate.Struct(req); err != nil {
		handleValidationErrors(w, err)
		return
	}

	// no need to catch the errors because we've already validated the data
	startTime, _ := time.Parse("15:04", req.StartDate)
	cinemaId, _ := uuid.Parse(req.CinemaID)
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	showTime, err := cfg.DB.CreateShowTime(r.Context(), database.CreateShowTimeParams{
		StartTime:      startTime,
		Price:          int32(req.Price),
		MovieID:        movieId,
		CinemaID:       cinemaId,
		ExperienceType: req.ExperienceType,
		StartDate:      startDate,
		EndDate:        endDate,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create movie show time", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, showTime)

}
