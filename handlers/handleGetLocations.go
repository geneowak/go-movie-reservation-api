package handlers

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleGetLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := cfg.DB.GetLocations(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "No locations found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching locations", err)
		return
	}

	respondWithJSON(w, http.StatusOK, locations)
}
