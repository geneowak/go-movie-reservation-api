package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/geneowak/go-expense-tracker/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type updateCinemaRequest struct {
	Name            string        `json:"name" validate:"required,min=3"`
	ExperienceTypes []string      `json:"experience_types" validate:"required,dive,min=1"`
	SeatMap         types.SeatMap `json:"seat_map" validate:"required"`
}

func (cfg *ApiConfig) handleUpdateCinemaDetails(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("cinemaId")
	cinemaId, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid cinema id", err)
		return
	}

	var req updateCinemaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleJsonDecodeError(w, err)
		return
	}

	if err := cfg.Validate.Struct(req); err != nil {
		handleValidationErrors(w, err)
		return
	}

	updatedCinema, err := cfg.DB.UpdateCinemaDetails(r.Context(), database.UpdateCinemaDetailsParams{
		ID:              cinemaId,
		Name:            req.Name,
		ExperienceTypes: req.ExperienceTypes,
		SeatMap:         req.SeatMap,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Cinema not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error updating cinema details", err)
		return
	}

	respondWithJSON(w, http.StatusOK, updatedCinema)
}
