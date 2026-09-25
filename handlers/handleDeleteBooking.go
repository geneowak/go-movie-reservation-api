package handlers

import (
	"errors"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleDeleteBooking(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("reservationId")
	reservationId, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid reservation Id", err)
		return
	}
	userId, _ := GetUserIdFromContext(r.Context())
	reservation, err := cfg.DB.GetUserBookingById(r.Context(), database.GetUserBookingByIdParams{
		ID:     reservationId,
		UserID: userId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Booking not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching booking", err)
		return
	}

	err = cfg.DB.DeleteUserBooking(r.Context(), database.DeleteUserBookingParams{
		ID:     reservation.ID,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting booking", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
