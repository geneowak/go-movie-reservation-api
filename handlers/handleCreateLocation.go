package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/database"
)

type createLocationRequest struct {
	Name         string `json:"name" validate:"required,min=3"`
	Address      string `json:"address" validate:"required,min=5"`
	GoogleMapUrl string `json:"google_map_url" validate:"required,url"`
}

func (cfg *ApiConfig) handleCreateLocation(w http.ResponseWriter, r *http.Request) {
	var req createLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request body", err)
		return
	}

	if err := cfg.Validate.Struct(req); err != nil {
		handleValidationErrors(w, err)
		return
	}

	location, err := cfg.DB.CreateLocation(r.Context(), database.CreateLocationParams{
		Name:         req.Name,
		Address:      req.Address,
		GoogleMapUrl: req.GoogleMapUrl,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create location", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, location)
}
