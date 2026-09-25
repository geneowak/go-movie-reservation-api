package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

	// a booked cinema should not be edited because this could alter some already booked seats
	cinemaBookings, err := cfg.DB.GetCinemaOngoingBookings(r.Context(), database.GetCinemaOngoingBookingsParams{
		ID:      cinemaId,
		EndDate: time.Now(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Cinema not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error checking cinema bookings.", err)
		return
	}
	var bookings []json.RawMessage
	_ = json.Unmarshal(cinemaBookings.Bookings, &bookings)
	if len(bookings) > 0 {
		respondWithError(w, http.StatusBadRequest, "Cinema has on going show and so can not be edited at the moment", errors.New("Cinema is currently booked"))
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
