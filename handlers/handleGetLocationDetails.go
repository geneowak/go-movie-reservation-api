package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/geneowak/cinehold/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleGetLocationDetails(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("locationId")
	locationId, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid location id", err)
		return
	}

	locationDetails, err := cfg.DB.GetLocationDetails(r.Context(), locationId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Location not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching location details", err)
		return
	}

	type response struct {
		database.GetLocationDetailsRow
		Cinemas json.RawMessage `json:"cinemas"`
	}

	respondWithJSON(w, http.StatusOK, response{
		locationDetails,
		locationDetails.Cinemas,
	})
}
