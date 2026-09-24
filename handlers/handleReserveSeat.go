package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
)

type reserveSeatRequest struct {
	SeatNo     string `json:"seat_no" validate:"required,is-valid-seat"`
	ShowTimeId string `json:"show_time_id" validate:"required,uuid_rfc4122"`
}

func (cfg *ApiConfig) handleReserveSeat(w http.ResponseWriter, r *http.Request) {
	var req reserveSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if typeErr, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
			msg := fmt.Sprintf("Invalid field %s: Expected type %s, but got %s", typeErr.Field, typeErr.Type.String(), typeErr.Value)
			respondWithError(w, http.StatusUnprocessableEntity, msg, err)
			return
		}
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
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
		if strings.Contains(err.Error(), EmptyResultSet) {
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
	cinema := results.Cinema
	log.Println("got cinema: ", cinema)
	// TODO: Validate the seat number

	// don't expect this to have an error since this handler is wrapped with the auth middleware
	userId, _ := UserIdFromContext(r.Context())
	// then we'll reserve the seat
	seatReserved, err := cfg.DB.CreateReservation(r.Context(), database.CreateReservationParams{
		ShowTimeID: showTime.ID,
		UserID:     userId,
		SeatNo:     req.SeatNo,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error researving seat number", err)
		return
	}

	log.Println("seat reserved:", seatReserved)

	respondWithJSON(w, http.StatusCreated, seatReserved)
}
