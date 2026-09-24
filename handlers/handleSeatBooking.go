package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
)

type seatBookingRequest struct {
	SeatReservations []string `json:"seat_reservations" validate:"required,dive,uuid_rfc4122"`
}

func (cfg *ApiConfig) handleSeatBooking(w http.ResponseWriter, r *http.Request) {
	var req seatBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleJsonDecodeError(w, err)
		return
	}

	if err := cfg.Validate.Struct(req); err != nil {
		handleValidationErrors(w, err)
		return
	}

	// going to first do the basic booking and will come back and validate
	// TODO: validate that the reservation belongs to the user and that the status is not in booked
	response := []database.Reservation{}
	userId, _ := UserIdFromContext(r.Context())
	for _, stringId := range req.SeatReservations {
		// don't expect an error here because we validated that they are uuids
		reservationId, _ := uuid.Parse(stringId)
		reservation, err := cfg.DB.MarkReservationBooked(r.Context(), database.MarkReservationBookedParams{
			ID:     reservationId,
			UserID: userId,
		})
		// not going to handle errors yet
		if err != nil {
			// TODO: we'll need to be doing all this in a transaction so that when an error occurs we rollback everything
			respondWithError(w, http.StatusInternalServerError, "Error creating reservation", err)
			return
		} else {
			response = append(response, reservation)
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}
