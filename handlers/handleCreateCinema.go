package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/geneowak/cinehold/internal/database"
	"github.com/geneowak/cinehold/internal/types"
	"github.com/google/uuid"
)

type createCinemaRequest struct {
	Name            string        `json:"name" validate:"required,min=3"`
	LocationId      uuid.UUID     `json:"location_id" validate:"required,uuid_rfc4122"`
	ExperienceTypes []string      `json:"experience_types" validate:"required,dive,min=1"`
	SeatMap         types.SeatMap `json:"seat_map" validate:"required"`
}

func (cfg *ApiConfig) handleCreateCinema(w http.ResponseWriter, r *http.Request) {
	var req createCinemaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleJsonDecodeError(w, err)
		return
	}

	if err := cfg.Validate.Struct(req); err != nil {
		handleValidationErrors(w, err)
		return
	}

	cinema, err := cfg.DB.CreateCinema(r.Context(), database.CreateCinemaParams{
		Name:            req.Name,
		LocationID:      req.LocationId,
		ExperienceTypes: req.ExperienceTypes,
		SeatMap:         req.SeatMap,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create cinema in DB", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, cinema)
}
