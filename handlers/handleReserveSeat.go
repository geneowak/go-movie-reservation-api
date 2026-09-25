package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type reserveSeatRequest struct {
	Seat       string `json:"seat" validate:"required,is_valid_seat"`
	ShowTimeId string `json:"show_time_id" validate:"required,uuid_rfc4122"`
}

func (cfg *ApiConfig) handleReserveSeat(w http.ResponseWriter, r *http.Request) {
	var req reserveSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleJsonDecodeError(w, err)
		return
	}

	if err := cfg.Validate.Struct(req); err != nil {
		handleValidationErrors(w, err)
		return
	}

	// no need to check error because we have already validated that it is a uuid
	ShowTimeId, _ := uuid.Parse(req.ShowTimeId)

	// let's first ensure that the show time is valid and that it is not a past event
	results, err := cfg.DB.GetShowTimeDetails(r.Context(), ShowTimeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Show time not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching show time details", err)
		return
	}
	showTime := results.ShowTime
	if showTime.EndDate.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusBadRequest, "Show time has already ended", err)
		return
	}
	// we'll get the cinema of the show time and validate that the seat no exists
	if !validateSeatNo(req, results.Cinema) {
		throwAsValidationError(w, "seat", "Seat number does not exist in show time cinema")
		return
	}

	// don't expect this to have an error since this handler is wrapped with the auth middleware
	userId, _ := GetUserIdFromContext(r.Context())

	// next is to validate if the seat has already been reserved or booked
	existingBooking, err := cfg.DB.GetReservationBySeatNo(r.Context(), database.GetReservationBySeatNoParams{
		ShowTimeID: showTime.ID,
		SeatNo:     req.Seat,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, "Failed to validate reservation", err)
			return
		}
	} else {
		if existingBooking.Status == "booked" {
			throwAsValidationError(w, "seat", "Seat number has already been booked")
			return
		}
		// user is only allowed to reserve a seat for 10 mins
		if existingBooking.ReservedAt != nil && time.Since(*existingBooking.ReservedAt) < 10*time.Minute {
			throwAsValidationError(w, "seat", "Seat number is currently reserved.")
			return
		}

		// reaching here means that the seat was reserved before but has gone stale i.e 10 mins have already passed
		updatedReservation, err := cfg.DB.UpdateReservation(r.Context(), database.UpdateReservationParams{
			ShowTimeID: showTime.ID,
			UserID:     userId,
			SeatNo:     req.Seat,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error updating seat reservation", err)
			return
		}

		respondWithJSON(w, http.StatusCreated, updatedReservation)
		return
	}
	// then we'll reserve the seat
	seatReserved, err := cfg.DB.CreateReservation(r.Context(), database.CreateReservationParams{
		ShowTimeID: showTime.ID,
		UserID:     userId,
		SeatNo:     req.Seat,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error researving seat number", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, seatReserved)
}

func validateSeatNo(req reserveSeatRequest, cinema database.Cinema) bool {
	// we have already validated that it must return only 2 items
	vals := strings.Split(req.Seat, ":")
	seatRow := vals[0]
	// this has already been validated so we are sure it will yield an int
	seatNo, _ := strconv.Atoi(vals[1])
	for _, row := range cinema.SeatMap.Rows {
		// if rows don't match then we continue
		if !strings.EqualFold(row.Row, seatRow) {
			continue
		}
		for _, seats := range row.Seats {
			if seats.Number == seatNo {
				// got a match so we exit early
				return true
			}
		}
		// row matched but seat was wrong so we exit from here
		return false
	}

	return false
}
