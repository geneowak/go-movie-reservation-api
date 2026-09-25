package handlers

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleGetBookings(w http.ResponseWriter, r *http.Request) {
	userId, _ := GetUserIdFromContext(r.Context())

	bookings, err := cfg.DB.GetUserBookings(r.Context(), userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithJSON(w, http.StatusNotFound, "No bookings at the moment")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching bookings", err)
		return
	}

	respondWithJSON(w, http.StatusOK, bookings)
}
