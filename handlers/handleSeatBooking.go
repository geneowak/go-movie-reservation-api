package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	userId, _ := GetUserIdFromContext(r.Context())

	errorMsgs := map[string][]string{}
	for i, id := range req.SeatReservations {
		// no error is expected here because we validated that they are uuids
		bookingId, _ := uuid.Parse(id)
		// bookings can only be made on non booked reservations that the user owns
		_, err := cfg.DB.GetUserBookingById(r.Context(), database.GetUserBookingByIdParams{
			ID:     bookingId,
			Status: "available",
			UserID: userId,
		})
		if err != nil {
			field := fmt.Sprintf("seat_reservations[%v]", i)
			if errors.Is(err, pgx.ErrNoRows) {
				errorMsgs[field] = append(errorMsgs[field], "Seat Reservation not found.")
			} else {
				errorMsgs[field] = append(errorMsgs[field], "Error checking seat reservation.")
			}
		}
	}
	if len(errorMsgs) > 0 {
		respondWithValidationErrors(w, errorMsgs)
		return
	}

	response := []database.Reservation{}
	for _, stringId := range req.SeatReservations {
		reservationId, _ := uuid.Parse(stringId)
		reservation, err := cfg.DB.MarkReservationBooked(r.Context(), database.MarkReservationBookedParams{
			ID:     reservationId,
			UserID: userId,
		})
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
